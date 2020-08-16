const admin = require("firebase-admin")
const five = require("johnny-five")
const board = new five.Board()

let RED
let YELLOW
let GREEN
let SERVO

board.on("ready", () => {
  console.log("board is ready!")
  RED = new five.Led(8)
  YELLOW = new five.Led(9)
  GREEN = new five.Led(10)
  SERVO = new five.Servo(11)

  RED.blink()
  YELLOW.blink()
  GREEN.blink()

  setTimeout(function () {
    RED.stop().off()
    YELLOW.stop().off()
    GREEN.stop().off()
    console.log("program is ready!")
  }, 2000)
})

admin.initializeApp({
  credential: admin.credential.cert("./key.json"),
  databaseURL: "https://chariot-rc.firebaseio.com",
})

const db = admin.database()
const ref = db.ref("chariot_1")

ref.on("child_changed", (snapshot) => {
  // FORWARD-BACKWARD MOVEMENT
  if (snapshot.key == "test_rotation") {
    rotation = snapshot.val()
    console.log(`UPDATE rotation, val: ${rotation}`)

    SERVO.to(rotation)
  }

  if (snapshot.key == "engine") {
    engine = snapshot.val()
    console.log(`UPDATE engine, val: ${engine}`)
    if (engine == 1) {
      SERVO.to(0)
      RED.off()
      GREEN.on()
      YELLOW.off()
    } else if (engine == 0) {
      SERVO.to(90)
      RED.off()
      GREEN.off()
      YELLOW.on()
    } else if (engine == -1) {
      SERVO.to(155)
      RED.on()
      GREEN.off()
      YELLOW.off()
    }
  }
})

process.on("exit", () => {
  RED.off()
  GREEN.off()
  YELLOW.off()
})
