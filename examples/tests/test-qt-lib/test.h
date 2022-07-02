#pragma once

#include <QtTest/QTest>

class MyQtLibTester : public QObject
{
    Q_OBJECT

public:
    MyQtLibTester() = default;

private:
    Q_SLOT void example_test_1();
    Q_SLOT void example_test_2();
};
