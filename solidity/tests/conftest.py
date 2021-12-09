#!/usr/bin/python3

import pytest

import brownie

from brownie import web3, TestLogicContract, SimpleLogicBatchMiddleware, Gravity, TestERC20A, ReentrantERC20, TestTokenBatchMiddleware, HashingTest, Contract

from eth_abi import encode_abi

from eth_account.messages import encode_defunct

@pytest.fixture(scope="session")
def signers(accounts):
    privKeys = getPrivKeys()
    acc_len = len(accounts)
    add_len = len(privKeys)
    if acc_len > add_len:
        return accounts[acc_len - add_len:]

    for i in range(0, add_len):
        accounts.add(privKeys[i])
        accounts[0].transfer(accounts[acc_len + i], 3466666 * 10 ** 18)
    return accounts[acc_len:]

def getPrivKeys():
    return [
        "0xc5e8f61d1ab959b397eecc0a37a6517b8e67a0e7cf1f4bce5591f3ed80199122",
        "0xd49743deccbccc5dc7baa8e69e5be03298da8688a15dd202e20f15d5e0e9a9fb",
        "0x23c601ae397441f3ef6f1075dcb0031ff17fb079837beadaf3c84d96c6f3e569",
        "0xee9d129c1997549ee09c0757af5939b2483d80ad649a0eda68e8b0357ad11131",
        "0x87630b2d1de0fbd5044eb6891b3d9d98c34c8d310c852f98550ba774480e47cc",
        "0x275cc4a2bfd4f612625204a20a2280ab53a6da2d14860c47a9f5affe58ad86d4",
        "0x7f307c41137d1ed409f0a7b028f6c7596f12734b1d289b58099b99d60a96efff",
        "0x2a8aede924268f84156a00761de73998dac7bf703408754b776ff3f873bcec60",
        "0x8b24fd94f1ce869d81a34b95351e7f97b2cd88a891d5c00abc33d0ec9501902e",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b29085",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b29086",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b29087",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b29088",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b29089",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908a",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908b",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908c",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908d",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908e",
        "0x28d1bfbbafe9d1d4f5a11c3c16ab6bf9084de48d99fbac4058bdfa3c80b2908f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf00",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf01",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf02",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf03",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf04",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf05",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf06",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf07",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf08",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf09",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0d",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0e",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf0f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf10",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf11",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf12",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf13",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf14",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf15",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf16",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf17",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf18",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf19",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1d",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1e",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf1f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf20",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf21",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf22",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf23",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf24",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf25",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf26",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf27",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf28",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf29",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2d",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2e",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf2f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf30",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf31",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf32",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf33",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf34",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf35",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf36",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf37",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf38",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf39",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3d",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3e",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf3f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf40",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf41",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf42",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf43",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf44",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf45",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf46",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf47",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf48",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf49",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4d",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4e",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf4f",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf50",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf51",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf52",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf53",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf54",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf55",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf56",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf57",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf58",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf59",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf5a",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf5b",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf5c",
        "0x2c7dd57db9fda0ea1a1428dcaa4bec1ff7c3bd7d1a88504754e0134b77badf5d",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb100",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb101",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb102",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb103",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb104",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb105",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb106",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb107",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb108",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb109",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10a",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10b",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10c",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10d",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10e",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb10f",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb110",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb111",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb112",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb113",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb114",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb115",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb116",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb117",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb118",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb119",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11a",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11b",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11c",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11d",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11e",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb11f",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb120",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb121",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb122",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb123",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb124",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb125",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb126",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb127",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb128",
        "0x47aa5fbb74b21f263888dfc24a7a7b184634142935d4e2152b1c901516eeb129",
        "0xb1bab011e03a9862664706fc3bbaa1b16651528e5f0e7fbfcbfdd8be302a13e7"
    ]

def getSignerAddresses(signers):
    ret = []
    for i in range(len(signers)):
        ret.append(signers[i].address)
    return ret

def makeCheckpoint(validators, powers, valsetNonce, gravityId):
    methodName = b"checkpoint"
    abiEncoded = encode_abi(["bytes32", "bytes32", "uint256", "address[]", "uint256[]"], [gravityId, methodName, valsetNonce, validators, powers])
    checkpoint = web3.keccak(abiEncoded)
    return checkpoint

def examplePowers():
    return [
        707,
        621,
        608,
        439,
        412,
        407,
        319,
        312,
        311,
        303,
        246,
        241,
        224,
        213,
        194,
        175,
        173,
        170,
        154,
        149,
        139,
        123,
        119,
        113,
        110,
        107,
        105,
        104,
        92,
        90,
        88,
        88,
        88,
        85,
        85,
        84,
        82,
        70,
        67,
        64,
        59,
        58,
        56,
        55,
        52,
        52,
        52,
        50,
        49,
        44,
        42,
        40,
        39,
        38,
        37,
        37,
        36,
        35,
        34,
        33,
        33,
        33,
        32,
        31,
        30,
        30,
        29,
        28,
        27,
        26,
        25,
        24,
        23,
        23,
        22,
        22,
        22,
        21,
        21,
        20,
        19,
        18,
        17,
        16,
        14,
        14,
        13,
        13,
        11,
        10,
        10,
        10,
        10,
        10,
        9,
        8,
        8,
        7,
        7,
        7,
        6,
        6,
        5,
        5,
        5,
        5,
        5,
        5,
        4,
        4,
        3,
        2,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1,
        1
    ]

def deployContracts(signers, gravityId, validators, powers, powerThreshold):
    testERC20 = TestERC20A.deploy({"from": signers[0]})
    valAddresses = getSignerAddresses(validators)
    checkpoint = makeCheckpoint(valAddresses, powers, 0, gravityId)

    GravityContract = web3.eth.contract(abi=Gravity.abi, bytecode=Gravity.bytecode)

    try:
        gas = GravityContract.constructor(gravityId, powerThreshold, valAddresses, powers).estimateGas({"from": signers[0].address})
    except ValueError as err:
        raise ValueError(err.args[0]["message"][50:])
    except brownie.exceptions.VirtualMachineError as err:
        raise ValueError(err.revert_msg)
    except BaseException as err:
        print(f"Unexpected {err=}, {type(err)=}")

    gravity = Gravity.deploy(gravityId, powerThreshold, valAddresses, powers, {"from": signers[0]})
    return gravity, testERC20, checkpoint

def signHash(signers, hash):
    sign = []
    for i in range(len(signers)):
        signed_message = web3.eth.account.sign_message(encode_defunct(hash), signers[i].private_key)
        sign.append([signed_message.v, signed_message.r, signed_message.s])
    return sign

def bstring2bytes32(str):
    return encode_abi(["bytes32"], [str])
