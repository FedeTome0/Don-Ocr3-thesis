// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract OracleQueue {
    uint256 public queryFee;
    uint256 public oracleReward;

    constructor(uint256 _queryFee, uint256 _oracleReward){
        require(_oracleReward <= _queryFee, "The reward cannot exceed the fee");
        queryFee = _queryFee;
        oracleReward = _oracleReward;
        modelCreator = msg.sender;
    }

    // The model creator is the owner of the contract
    address public modelCreator;

    // Separated counters for the two mappings
    uint256 public requestCounter;
    uint256 public oracleJobCounter;

    // Structure for the Cosutomer queue
    struct CustomerRequest {
        string ipfsCid;
        address requester;
        uint256 payment;
        bool isProcessed; // To track if the model creator has already processed it
    }

    struct OracleJob {
        uint256 originalRequestId;
        string ipfsCid;
    }

    // Mapping to keep track of the jobs
    mapping(uint256 => CustomerRequest) public customerQueue;
    mapping(uint256 => OracleJob) public oracleQueue;

    // Event 1: Emitted when the Customer pays and asks for a model computation
    // Listened off-chain by the modelCreator
    event LogNewCustomerRequest(
        uint256 indexed requestId, 
        string ipfsCid, 
        address requester, 
        uint256 payment
    );

    // Event 2: Emitted when the Model Creator approves the job
    // Listened off-chain by the Oracle Network
    event LogNewJobForOracles(
        uint256 indexed jobId, 
        string ipfsCid
    );

    modifier onlyOwner() {
        require(msg.sender == modelCreator, "Only Model Creator can do this");
        _;
    }

    // =========================================================
    // PHASE 1: CUSTOMER
    // =========================================================    
    function requestAttribution(string calldata _ipfsCid) external payable {
        // Base check on the exact payment
        require(msg.value == queryFee, "Amount error: must pay the rigth queryFee"); 

        uint256 currentReqId = requestCounter;

        // Save the new request in the queue
        customerQueue[currentReqId] = CustomerRequest({
            ipfsCid: _ipfsCid,
            requester: msg.sender,
            payment: msg.value,
            isProcessed: false
        });

        emit LogNewCustomerRequest(currentReqId, _ipfsCid, msg.sender, msg.value);
        requestCounter++;
    }

    // =========================================================
    // PHASE 2: MODEL CREATOR
    // =========================================================
    // The Model Creator calls this function after validating the request
    function approveJob(uint256 _requestId) external onlyOwner {
        // Fetch from the customer queue
        CustomerRequest storage req = customerQueue[_requestId];
        require(!req.isProcessed, "Request already approved");

        // Update the state
        req.isProcessed = true;

        // Prepare new ID for the Oracle queue
        uint256 currentOracleJobId = oracleJobCounter;
        // Insert the approved job into the oracleQueue
        oracleQueue[currentOracleJobId] = OracleJob({
            originalRequestId: _requestId,
            ipfsCid: req.ipfsCid
        });

        // Emit the event that wakes up the Oracles network
        emit LogNewJobForOracles(currentOracleJobId, req.ipfsCid);

        oracleJobCounter++;
    }

    // =========================================================
    // PHASE 3: PAYMENT
    // =========================================================
    address public oracleVerifierAddress;

    // The model creator connects the queue to the Verifier after deployment
    function setVerifierAddress(address _verifier) external onlyOwner {
        oracleVerifierAddress = _verifier;
    }

    // This function is called automatically by the OracleVerifier contract
    // at the end of "transmit" function to refund the Oracle that executed the job
    function rewardOracle(address payable _oracle) external {
        require(msg.sender == oracleVerifierAddress, "Only the Verifier can unlock the funds");

        // If there is a fee, reimburse the Oracle for the spent gas
        if (oracleReward > 0) {
            (bool success, ) = _oracle.call{value: oracleReward}("");
            require(success, "Refund to the oracle failed");
        }
    } 
    
}
