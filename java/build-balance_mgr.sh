CUR_PATH=${PWD}
#balance_mgr
#maven
cd balance_mgr/maven
rm -rf target
mvn package
#gradle
cd ../gradle
rm -rf build
gradle task shadowJar

