#include <QtCore/QDebug>

#include <cc-lib/mylib.h>

int main(int argc, char** argv)
{
    qDebug() << "Hello:" << mylib::add(1, 2);

    return 0;
}
