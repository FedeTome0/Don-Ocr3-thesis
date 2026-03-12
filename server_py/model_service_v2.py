import os
import sys
import uuid
import time
import json
import random
import threading
import subprocess
import configparser
import queue
import numpy as np
import torch
from contextlib import nullcontext
from flask import Flask, request, jsonify

# --- MODEL IMPORTS ---
import model_utils
from dattri.benchmark.models.nanoGPT.model import GPT, GPTConfig

# --- CONFIGURATION ---
CONFIG_FILE = 'model_service_config.ini'
config = configparser.ConfigParser()
config.read(CONFIG_FILE)

APP_PORT = int(config['APPLICATION']['port'])
APP_HOST = config['APPLICATION']['host']

# Model parameters
meta_path = config['MODEL']['meta_path']
checkpoint_path = config['MODEL']['checkpoint_path']
device = config['MODEL']['device']
block_size = int(config['MODEL']['block_size'])
seed = int(config['MODEL']['seed'])
max_new_tokens = int(config['MODEL']['max_new_tokens'])
temperature = float(config['MODEL']['temperature'])
top_k = int(config['MODEL']['top_k'])

TEMP_DIR = "./tmp_jobs"
os.makedirs(TEMP_DIR, exist_ok=True)

app = Flask(__name__)

# SHARED MEMORY AND LOCK
JOBS = {}
jobs_lock = threading.Lock()
JOB_QUEUE = queue.Queue()

# Global model variables
ctx = nullcontext()
torch.manual_seed(seed)
model = None
encode = None
decode = None

def load_model_logic():
    global model, encode, decode
    print("Loading model...")
    encode_f, decode_f = model_utils.load_meta(meta_path)
    encode = encode_f
    decode = decode_f

    checkpoint = torch.load(checkpoint_path, map_location=device)
    gptconf = GPTConfig(**checkpoint['model_args'])
    model = GPT(gptconf)
    state_dict = checkpoint['model']

    unwanted_prefix = '_orig_mod.'
    for k,v in list(state_dict.items()):
        if k.startswith(unwanted_prefix):
            state_dict[k[len(unwanted_prefix):]] = state_dict.pop(k)

    model.load_state_dict(state_dict)
    model.eval()
    model.to(device)
    print("Model loaded and ready.")

def generate_text(prompt):
    start_ids = encode(prompt)
    x = (torch.tensor(start_ids, dtype=torch.long, device=device)[None, ...])
    with torch.no_grad():
        with ctx:
            y = model.generate(x, max_new_tokens, temperature=temperature, top_k=top_k)
            return decode(y[0].tolist())

def run_full_process(job_id, prompt):
    print(f"[JOB {job_id}] Step 1: Generating text...")

    input_filename = os.path.join(TEMP_DIR, f"{job_id}.in")
    output_filename = os.path.join(TEMP_DIR, f"{job_id}.out")

    try:
        # A. Generation
        generated_story = generate_text(prompt)
        print(f"[JOB {job_id}] Generated text (len={len(generated_story)}).")

        # B. Write to file
        with open(input_filename, "w", encoding="utf-8") as f:
            f.write(generated_story)

        # C. Attribution (Subprocess)
        print(f"[JOB {job_id}] Step 2: Starting Attribution (Subprocess)...")
        cmd = [sys.executable, "model_attributor.py", input_filename, output_filename]

        # Run external script and capture output/errors
        result = subprocess.run(cmd, capture_output=True, text=True)

        if result.returncode != 0:
            # If it fails, print the stderr from the subprocess
            print(f"STDERR from subprocess:\n{result.stderr}")
            raise Exception(f"Attributor Error: {result.stderr}")

        if not os.path.exists(output_filename):
            raise Exception("Output file not found (Attributor failed silently?)")

        # D. Read and Aggregate results
        print(f"[JOB {job_id}] Step 3: Reading results...")
        raw_data = np.loadtxt(output_filename)

        if raw_data.ndim == 2:
            aggregated_scores = np.mean(raw_data, axis=1)
        elif raw_data.ndim == 1:
            aggregated_scores = raw_data
        else:
            aggregated_scores = np.array([raw_data])

        if aggregated_scores.ndim == 0: aggregated_scores = np.array([aggregated_scores])

        # E. Convert to BigInt strings
        blockchain_ready_scores = []
        for score in aggregated_scores:
            # Protection against NaN or Infinity
            if np.isnan(score) or np.isinf(score): score = 0.0

            int_val = int(float(score) * 1e18)
            blockchain_ready_scores.append(str(int_val))

        # F. Save Result
        with jobs_lock:
            JOBS[job_id]["result"] = blockchain_ready_scores
            JOBS[job_id]["status"] = "completed"

        print(f"[JOB {job_id}] COMPLETED SUCCESSFULLY!")

    except Exception as e:
        print(f"[JOB {job_id}] CRITICAL ERROR: {str(e)}")
        with jobs_lock:
            JOBS[job_id]["status"] = "error"
            JOBS[job_id]["error"] = str(e)

    finally:
        # Clean up input file, keep output for debug if needed
        if os.path.exists(input_filename): os.remove(input_filename)

# --- BACKGROUND WORKER (SEQUENTIAL EXECUTION) ---
def background_worker():
    """
    This thread runs in background and takes only one job at a time
    from the queue and execute the AI. This grants no parallel executions
    """
    while True:
        job_id, prompt = JOB_QUEUE.get()

        # Update the status to "processing" only when the job starts
        with jobs_lock:
            if job_id in JOBS:
                JOBS[job_id]["status"] = "processing"

        try:
            run_full_process(job_id, prompt)
        except Exception as e:
            print(f"[WORKER] Unexpected error for the job {job_id}: {e}")
        finally:
            # Alert the queue that this task is done, unlocking the next
            JOB_QUEUE.task_done()

# --- ENDPOINTS ---

@app.route('/attribute', methods=['POST'])
def attribute():
    data = request.json
    if not data or 'text' not in data:
        return jsonify({"error": "Missing 'text' field"}), 400

    # --- ROBUST ID LOGIC ---
    # Reads job_id (new standard) OR cid (old standard) OR creates one
    job_id = data.get('job_id') or data.get('cid') or str(uuid.uuid4())

    prompt = data['text']

    with jobs_lock:
        if job_id in JOBS:
            print(f"--> [DEDUPLICATION] Duplicate request for Job {job_id}. Ignoring.")
            return jsonify({
                "message": "Job already exists",
                "job_id": job_id,
                "status": JOBS[job_id]["status"]
            }), 200

        JOBS[job_id] = {"status": "queued", "result": None}

    print(f"--> [NEW] Queuing Job {job_id}")
    # QUEUE THE JOB (DO NOT START A THREAD HERE)
    JOB_QUEUE.put((job_id, prompt))

    return jsonify({"message": "Job Queued", "job_id": job_id}), 202

@app.route('/result/<job_id>', methods=['GET'])
def get_result(job_id):
    job = JOBS.get(job_id)
    if not job:
        return jsonify({"error": "Job not found"}), 404
        
    return jsonify(job), 200

if __name__ == '__main__':
    load_model_logic()

    # Start the background worker before exposing the API
    threading.Thread(target=background_worker, daemon=True).start()

    print(f"Starting server on {APP_HOST}:{APP_PORT}...")
    # debug=False to avoid double loading or thread issues
    app.run(host=APP_HOST, port=APP_PORT, threaded=True, debug=False)