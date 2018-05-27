import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQuick.Controls 1.4
import QtQuick.Window 2.3

ApplicationWindow {
    id:window
    title: "BtcQuant"
    width: Screen.desktopAvailableWidth /2
	height: Screen.desktopAvailableHeight /2
    minimumWidth: 500
	minimumHeight: 400
    // color: "gray"

    TabView {
        id :windowTab

        anchors.centerIn: parent
        width: parent.width
        height: parent.height
        
        Tab {
            title: "行情"
            Rectangle {
                width: parent.width
                height: parent.height
                TableView {
                    id: tableView1
                    anchors.top: parent.top
                    anchors.left: parent.left
                    width: parent.width*2/5
                    height: parent.height
                    TableViewColumn {
                        role: "symbol1"
                        title: "Symbol1"
                    }
                    TableViewColumn {
                        role: "test1"
                        title: "Test1"
                    }
                }

                TableView {
                    id: tableView2
                    anchors.top: parent.top
                    anchors.right: parent.right
                    width: parent.width*3/5
                    height: parent.height
                    TableViewColumn {
                        role: "symbol2"
                        title: "Symbol2"
                    }
                    TableViewColumn {
                        role: "test2"
                        title: "Test2"
                    }
                }
            }
            
        }

        Tab {
            title: "待开发"
            Rectangle {color: "gray"}
        }
    }
}