'use strict';

var Node = require('../../../Node');
var BigNumber = require('bignumber.js');
var expect = require('chai').expect;
var sleep = require("system-sleep")
var os = require('os');

var count = 6;
var validators = 6;
var blockInterval = 5;
var dynastyInterval = 60;
var reward = new BigNumber("48e16");
var initial = new BigNumber("1e18");
var nodes = new Node(count);
var now = new Date().getTime() / 1000;
nodes.Start();

describe('seed server start correctly', function () {
    before(function (done) {
        this.timeout(10000);
        setTimeout(done, 4000);
    });

    it('check status', function () {
        while (true) {
            var nodeInfo = nodes.RPC(0).api.nodeInfo()
            console.log(nodeInfo)
            if (nodeInfo.route_table.length == count - 1) {
                break;
            }
            sleep(3000);
        }
    });
});

describe('change candidates & votes', function () {
    it('change candidates', function () {
        var node = nodes.Node(1);
        var resp = node.RPC().admin.unlockAccount(node.Coinbase(), node.Passphrase());
        console.log(resp);
        expect(resp).to.not.have.property("error");
        resp = node.RPC().api.sendTransaction(
            node.Coinbase(), node.Coinbase(),
            0, 1, 0, 20000000,
            null, { action: "logout" }, null);
        console.log(resp);
        expect(resp).to.be.have.property("txhash");
    });

    it('change votes', function () {
        var node = nodes.Node(1);
        var resp = node.RPC().admin.unlockAccount(node.Coinbase(), node.Passphrase());
        console.log(resp);
        expect(resp).to.not.have.property("error");
        resp = node.RPC().api.sendTransaction(
            node.Coinbase(), node.Coinbase(),
            0, 2, 0, 20000000,
            null, null, { action: "do", delegatee: "fc751b484bd5296f8d267a8537d33f25a848f7f7af8cfcf6" });
        console.log(resp);
        expect(resp).to.be.have.property("txhash");
    });
});

describe('check dynasty', function () {
    it('check dynasty', function () {
        var node = nodes.Node(0);
        for (var i = 0; i < 100; i++) {
            var dynasty = parseInt((new Date().getTime() / 1000 - now) / dynastyInterval);
            var resp = node.RPC().admin.getDynasty();
            console.log(dynasty);
            console.log(resp);
            if (dynasty > 3) {
                expect(resp.delegatees[1]).to.be.equal("333cb3ed8c417971845382ede3cf67a0a96270c05fe2f700");
                expect(resp.delegatees[validators - 1]).to.be.equal("fc751b484bd5296f8d267a8537d33f25a848f7f7af8cfcf6");
                break;
            }
            sleep(5000);
        }
    });
});

describe('quit', function () {
    it('quit', function () {
        nodes.Stop();
    });
});
