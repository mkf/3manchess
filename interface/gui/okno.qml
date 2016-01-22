//import GoExtensions 1.0
import QtQuick 2.4
import QtQuick.Window 2.2

Window {
	id: win

	width: 1000
	height: 700
	title: "3manchess"
	visible: true

	/*Rectangle {
		width: 699; height: 699; color: "black"

		Rectangle {
			width: 698; height: 698; color: "white"; border.color: "gray"; border.width: 1
			radius: 349
		}
	}*/

	Image {
		source: "3manchesstemp700px.png"
		MouseArea {
			anchors.fill: parent
			onClicked: console.log("PLACEHOLDER")
		}
	}
}
