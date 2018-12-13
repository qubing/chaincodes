package chaincode.example;

import org.hyperledger.fabric.shim.ChaincodeBase;
import org.hyperledger.fabric.shim.ChaincodeStub;

import java.util.List;

public class BalanceManager extends ChaincodeBase {
    @Override
    public Response init(ChaincodeStub stub) {
        try {
            String func = stub.getFunction();
            if (func == null || func.trim().equals("") || func.equals("init")) {
                return _init(stub);
            } else if (func.equals("deploy")) {
                return _deploy(stub);
            }
            throw new Exception("function name not correct. expected: 'init' or 'upgrade', actual: '" + func + "'");
        } catch (Throwable e) {
            return newErrorResponse(e);
        }
    }

    private Response _init(ChaincodeStub stub) throws Throwable {
        return newSuccessResponse();
    }

    private Response _deploy(ChaincodeStub stub) throws Throwable {
        return newSuccessResponse();
    }

    @Override
    public Response invoke(ChaincodeStub stub) {
        try {
            String func = stub.getFunction();
            if (func == null || func.trim().equals("")) {
                throw new Exception("function name not provided.");
            }

            if (func.equals("create")) {
                return _create(stub);
            }

            if (func.equals("charge")) {
                return _charge(stub);
            }

            if (func.equals("query")) {
                return _query(stub);
            }

            throw new Exception("function name not correct. expected: 'create', 'charge' or 'query', actual: '" + func + "'");
        } catch (Throwable e) {
            return newErrorResponse(e);
        }
    }

    private Response _create(ChaincodeStub stub) throws Throwable {
        List<String> args = stub.getParameters();
        if (args.size() != 1 && args.size() != 2) {
            return newErrorResponse("Incorrect number of arguments. Expecting 1 or 2");
        }
        String account = args.get(0);
        String initBalance = "0";
        if (args.size() == 2) {
            initBalance = args.get(1);
            try {
                Integer.parseInt(initBalance);
            } catch (NumberFormatException e) {
                throw new Exception("init balance not valid.");
            }
        }

        stub.putStringState(account, initBalance);

        return newSuccessResponse();
    }

    private Response _charge(ChaincodeStub stub) throws Throwable {
        List<String> args = stub.getParameters();
        if (args.size() != 2) {
            return newErrorResponse("Incorrect number of arguments. Expecting 2");
        }
        String account = args.get(0);
        String charge = args.get(1);
        int chargeAmount = 0;
        int balanceAmount = 0;
        try {
            chargeAmount = Integer.parseInt(charge);
        } catch (NumberFormatException e) {
            throw new Exception(String.format("Error: amount to charge should be number. actual: %s", charge));
        }

        String balance = stub.getStringState(account);

        if (balance == null) {
            throw new Exception(String.format("Error: account %s to charge not existing.", account));
        }
        try {
            balanceAmount = Integer.parseInt(balance);
        } catch (NumberFormatException e) {
            throw new Exception(String.format("Error: amount to charge should be number. actual: %s", charge));
        }

        balanceAmount += chargeAmount;

        stub.putStringState(account, String.valueOf(balanceAmount));

        return newSuccessResponse();
    }

    private Response _query(ChaincodeStub stub) throws Throwable {
        List<String> args = stub.getParameters();
        if (args.size() != 1) {
            return newErrorResponse("Incorrect number of arguments. Expecting 1");
        }
        String account = args.get(0);

        String balance	= stub.getStringState(account);
        if (balance == null) {
            return newErrorResponse(String.format("Error: balance for %s is null", account));
        }

        return newSuccessResponse(balance/*, ByteString.copyFrom(balance, UTF_8).toByteArray()*/);
    }


    public static void main(String[] args) {
        new BalanceManager().start(args);
    }
}
