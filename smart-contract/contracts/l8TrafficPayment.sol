pragma solidity ^0.8.28;

error TransferFailed();

contract L8TrafficPayment {
    address public immutable owner;
    address payable public receiver;

    event TrafficPaid(string clientID, address payer, uint amount);

    constructor(address payable _receiver) {
        owner = msg.sender;
        receiver = _receiver;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner is allowed to execute this function");
        _;
    }

    function setReceiverAddress(address payable _newReceiverAddress) public onlyOwner {
        receiver = _newReceiverAddress;
    }

    function pay(string calldata clientID) external payable {
        (bool success, ) = receiver.call{value: msg.value}("");
        if (!success) {
            revert TransferFailed();
        }

        emit TrafficPaid(clientID, msg.sender, msg.value);
    }
}