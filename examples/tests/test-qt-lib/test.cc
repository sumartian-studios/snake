#include "test.h"

void MyQtLibTester::example_test_1()
{
    QVERIFY(1 != 2);
}

void MyQtLibTester::example_test_2()
{
    QVERIFY(1 != 2);
}

QTEST_MAIN(MyQtLibTester)
