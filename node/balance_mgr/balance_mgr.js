const shim = require('fabric-shim');
const util = require('util');

var BalanceManager = class {

  // Initialize the chaincode
  async Init(stub) {
    console.info('========= balance_mgr Init =========');
    let ret = stub.getFunctionAndParameters();
    console.info(ret);

    let method = this[ret.fcn];
    if (!method) {
      console.log('no method of name:' + ret.fcn + ' found');
      return shim.success();
    }

    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  async Invoke(stub) {
    let ret = stub.getFunctionAndParameters();
    console.info(ret);
    let method = this[ret.fcn];
    if (!method) {
      console.log('no method of name:' + ret.fcn + ' found');
      return shim.success();
    }
    try {
      let payload = await method(stub, ret.params);
      return shim.success(payload);
    } catch (err) {
      console.log(err);
      return shim.error(err);
    }
  }

  async init(stub, args) {
    //to be done
  }

  async upgrade(stub, args) {
    //to be done
  }

  async create(stub, args) {
    if (args.length != 1) {
      throw new Error('Incorrect number of arguments. Expecting 1');
    }

    let account = args[0];
    if (!account) {
      throw new Error('asset holding must not be empty');
    }

    // Get the state from the ledger
    let valbytes = await stub.getState(account);
    if (valbytes) {
      throw new Error('asset holder :' + account + ' already existing');
    }

    // Write the states back to the ledger
    await stub.putState(account, Buffer.from("0"));
  }

  async charge(stub, args) {
    if (args.length != 2) {
      throw new Error('Incorrect number of arguments. Expecting 2');
    }

    let account = args[0];
    if (!account) {
      throw new Error('asset holding must not be empty');
    }

    let val = args[1];
    if (!val || typeof parseInt(val) !== 'number') {
      throw new Error('charge amount must not be empty or non-number');
    }

    // Get the state from the ledger
    let balbytes = await stub.getState(account);
    if (!balbytes) {
      throw new Error('asset holder :' + account + " doesn't existing");
    }

    let bal = parseInt(balbytes.toString());

    if (typeof bal !== 'number') {
      bal = parseInt(val);
    } else {
      bal += parseInt(val);
    }

    // Write the states back to the ledger
    await stub.putState(account, Buffer.from(bal.toString()));
  }

  // query callback representing the query of a chaincode
  async query(stub, args) {
    if (args.length != 1) {
      throw new Error('Incorrect number of arguments. Expecting name of the person to query')
    }

    let jsonResp = {};
    let A = args[0];

    // Get the state from the ledger
    let Avalbytes = await stub.getState(A);
    if (!Avalbytes) {
      jsonResp.error = 'Failed to get state for ' + A;
      throw new Error(JSON.stringify(jsonResp));
    }

    jsonResp.name = A;
    jsonResp.amount = Avalbytes.toString();
    console.info('Query Response:');
    console.info(jsonResp);
    return Avalbytes;
  }
};

shim.start(new BalanceManager());