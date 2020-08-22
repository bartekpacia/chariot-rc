const admin = require("firebase-admin")
const five = require("johnny-five")
const Raspi = require("raspi-io").RaspiIO
const board = new five.Board({
  io: new Raspi(),
})

let RED
let YELLOW
let GREEN
let SERVO

board.on("ready", () => {
  const led = new five.Led("GPIO14")
  led.blink()
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
      SERVO.to(45)
      RED.off()
      GREEN.on()
      YELLOW.off()
    } else if (engine == 0) {
      SERVO.to(90)
      RED.off()
      GREEN.off()
      YELLOW.on()
    } else if (engine == -1) {
      SERVO.to(170)
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
