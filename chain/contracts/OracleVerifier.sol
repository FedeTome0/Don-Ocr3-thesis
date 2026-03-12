// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Add the oracleQueue interface
interface IOracleQueue {
    function rewardOracle(address payable _oracle) external;
}

contract OracleVerifier {

    IOracleQueue public oracleQueue;

    struct Result {
        int128[] flatMatrix;
        address  submitter;
        uint256  timestamp;
        bool     saved;
    }

    mapping(uint256 => Result) public results;
    mapping(uint64  => bool)   public usedSeqNr; // To avoid replay attacks

    mapping(address => bool) public isOracle;
    uint8 public f;
    bytes32 public expectedConfigDigest;

    event JobCompleted(uint256 indexed jobId, address indexed submitter, uint256 vectorLength, uint256 timestamp);

    // Constructor
    constructor(address[] memory _oracles, uint8 _f, bytes32 _configDigest, address _queueAddress) {
        require(_oracles.length >= 3 * _f + 1, "Network too small for the given f");
        for (uint i = 0; i < _oracles.length; i++) {
            isOracle[_oracles[i]] = true;
        }
        f = _f;
        expectedConfigDigest = _configDigest;

        // Connect the interface to the specified address
        oracleQueue = IOracleQueue(_queueAddress);
    }

    // =====================================================================
    // CUSTOM HASH FUNCTION (Perfect mirror of the Go keyring logic)
    // =====================================================================
    function _computeCustomHash(
        bytes32 configDigest,
        uint64 seqNr,
        bytes calldata report
    ) internal pure returns (bytes32) {
        // Usa encodePacked per concatenare i byte senza padding (uguale all'append in Go)
        return keccak256(abi.encodePacked(configDigest, seqNr, report));
    }

    // =========================================================================
    // TRANSMIT FUNCTION (Standard OCR signature verification)
    // =========================================================================    
    function transmit(
        bytes32 configDigest,
        uint64 seqNr,
        bytes calldata report,
        bytes32[] calldata rs, 
        bytes32[] calldata ss, 
        bytes32 rawVs // A single bytes32 parameter containing up to 32 packed 'v' values
    ) external {
        // Basic checks and Anti-Replay protection
        require(configDigest == expectedConfigDigest, "Invalid ConfigDigest");
        require(!usedSeqNr[seqNr], "Sequence Number already used (Replay Attack)");

        // Mark the sequence number as used
        usedSeqNr[seqNr] = true;

        // 1. Calculate the hash using raw byte concatenation
        bytes32 signedHash = _computeCustomHash(configDigest, seqNr, report);

        uint256 validSignatures = 0;

        // Create an in-memory array to track signers and prevent duplicates
        address[] memory seenSigners = new address[](rs.length);

        for (uint i = 0; i < rs.length; i++) {
            uint8 v = uint8(rawVs[i]);
            
            // 2. Perform ecrecover, retrieving the address of the oracles
            address signer = ecrecover(signedHash, v, rs[i], ss[i]);
            // check oracle signature
            require(isOracle[signer], "Unauthorized signature");

            // Check for duplicate signatures
            for (uint j = 0; j < i; j++) {
                require(seenSigners[j] != signer, "Duplicate signature found");
            }
            seenSigners[i] = signer;
            validSignatures++;
        }

        // Check the required signatures
        require(validSignatures >= f + 1, "BFT Quorum not reached");

        // Decoding and Storage
        (uint256 jobId, int128[] memory _flatMatrix) = abi.decode(report, (uint256, int128[]));
        // Revert if the job is already saved by another oracle
        require(!results[jobId].saved, "Job already completed");

        results[jobId] = Result({
            flatMatrix:     _flatMatrix,
            submitter:      msg.sender,
            timestamp:      block.timestamp,
            saved:          true
        });

        emit JobCompleted(jobId, msg.sender, _flatMatrix.length, block.timestamp);

        // unlock the funds and reimburse the oracle's gas costs
        oracleQueue.rewardOracle(payable(msg.sender));
    }

    function getResult(uint256 _jobId) external view returns (int128[] memory, address, uint256) {
        require(results[_jobId].saved, "Result not found");
        return (results[_jobId].flatMatrix, results[_jobId].submitter, results[_jobId].timestamp);
    }

    function isCompleted(uint256 _jobId) external view returns (bool) {
        return results[_jobId].saved;
    }
}
/*pragma solidity ^0.8.0;

contract OracleVerifier {
    // struct for keeping track of the matrixes
    struct Result {
        int128[] flatMatrix;
        address submitter;
        uint256 timestamp;
        bool saved;
    }

    // mapping for the history of the results
    // Grant access directly to the data without cycle. So that the cost is low
    mapping(uint256 => Result) public results;

    //string public latestOutcome;
    //uint256 public reportCount;
    // Mapping to remember which RequestID are already saved
    //mapping(bytes32 => bool) public processedRequests;

    event JobCompleted(uint256 indexed jobId, address indexed submitter, uint256 vectorLenght, uint256 timestamp);

    // This is the function called by the oracles in GO
    function saveOutcome(uint256 jobId, int128[] calldata _flatMatrix) external {
        require(!results[jobId].saved, "Request already fulfilled");

        // Save on-chain (storage)
        results[jobId] = Result({
            flatMatrix:     _flatMatrix,
            submitter:  msg.sender,
            timestamp:  block.timestamp,
            saved:     true
        });
        emit JobCompleted(jobId, msg.sender,_flatMatrix.length, block.timestamp);
    }

    // Getter to retrieve the matrix (for the smart licenses)
    function getResult(uint256 _jobId) external view returns (int128[] memory, address, uint256) {
        require(results[_jobId].saved, "Result not found");

        Result memory r = results[_jobId];
        return (r.flatMatrix, r.submitter, r.timestamp);
    }

    // Check if the job is already completed (for the transmitter)
    function isCompleted(uint256 _jobId) external view returns (bool) {
        return results[_jobId].saved;
    }

}
*/