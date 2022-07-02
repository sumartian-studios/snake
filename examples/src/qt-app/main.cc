#include <cc-lib/mylib.h>
#include <cstdlib>
#include <iostream>
#include <termcolor/termcolor.hpp>

#include <QtCore/QCoreApplication>
#include <QtCore/QDebug>
#include <QtCore/QFile>

int main(int argc, char** argv)
{
    QCoreApplication app(argc, argv);

    QFile file(":/data.txt");

    if (!file.open(QIODevice::ReadOnly | QIODevice::Text)) {
        qCritical() << "Could not read text file";
        return EXIT_FAILURE;
    }

    qDebug() << file.readAll();
    qDebug() << mylib::add(2, 4);

    std::cout << termcolor::red << "Hello, ";                // 16 colors
    std::cout << termcolor::color<100> << "Colorful ";       // 256 colors
    std::cout << termcolor::color<211, 54, 130> << "World!"; // true colors
    std::cout << std::endl;

    return EXIT_SUCCESS;
}
