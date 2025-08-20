// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReceiptRegistry {
    struct Receipt {
        bytes32 hash;
        string kid;
        uint256 timestamp;
        bool exists;
    }
    
    mapping(string => Receipt) public receipts;
    
    event ReceiptAnchored(string receiptId, bytes32 hash, string kid, uint256 timestamp);
    
    function anchorReceipt(string memory receiptId, bytes32 hash, string memory kid) public {
        require(!receipts[receiptId].exists, "Receipt already exists");
        
        receipts[receiptId] = Receipt({
            hash: hash,
            kid: kid,
            timestamp: block.timestamp,
            exists: true
        });
        
        emit ReceiptAnchored(receiptId, hash, kid, block.timestamp);
    }
    
    function verifyReceipt(string memory receiptId, bytes32 hash) public view returns (bool, string memory, uint256) {
        Receipt memory receipt = receipts[receiptId];
        require(receipt.exists, "Receipt does not exist");
        
        return (receipt.hash == hash, receipt.kid, receipt.timestamp);
    }
}
