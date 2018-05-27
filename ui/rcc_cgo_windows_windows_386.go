package main

/*
#cgo CFLAGS: -fno-keep-inline-dllexport -O2 -Wextra -Wall -W -DUNICODE -D_UNICODE -DQT_NEEDS_QMAIN -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -fno-keep-inline-dllexport -O2 -std=gnu++11 -Wextra -Wall -W -fexceptions -mthreads -DUNICODE -D_UNICODE -DQT_NEEDS_QMAIN -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB
#cgo CXXFLAGS: -I../../BtcQuant -I. -ID:/Qt/Qt5.10.0/5.10.0/mingw53_32/include -ID:/Qt/Qt5.10.0/5.10.0/mingw53_32/include/QtGui -ID:/Qt/Qt5.10.0/5.10.0/mingw53_32/include/QtANGLE -ID:/Qt/Qt5.10.0/5.10.0/mingw53_32/include/QtCore -Irelease -ID:/Qt/Qt5.10.0/5.10.0/mingw53_32/mkspecs/win32-g++
#cgo LDFLAGS:        -Wl,-s -Wl,-subsystem,windows -mthreads
#cgo LDFLAGS:        -lmingw32 -LD:/Qt/Qt5.10.0/5.10.0/mingw53_32/lib D:/Qt/Qt5.10.0/5.10.0/mingw53_32/lib/libqtmain.a -LC:/utils/my_sql/my_sql/lib -LC:/utils/postgresql/pgsql/lib -lshell32 D:/Qt/Qt5.10.0/5.10.0/mingw53_32/lib/libQt5Gui.a D:/Qt/Qt5.10.0/5.10.0/mingw53_32/lib/libQt5Core.a
#cgo LDFLAGS: -Wl,--allow-multiple-definition
#cgo CFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
#cgo CXXFLAGS: -Wno-unused-parameter -Wno-unused-variable -Wno-return-type
*/
import "C"
