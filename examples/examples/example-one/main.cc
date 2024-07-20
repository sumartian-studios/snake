#include <cc-lib/mylib.h>

#include <QtCore/QCoreApplication>
#include <QtCore/QDebug>

int main(int argc, char** argv)
{
    qDebug() << "Hello:" << mylib::add(1, 2);

    return 0;
}
