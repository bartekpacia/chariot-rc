const admin = require("firebase-admin")
const five = require("johnny-five")
const Raspi = require("raspi-io").RaspiIO
let board

try {
  board = new five.Board({
  	io: new Raspi()
  })	
} catch (err) {
  console.log("failed to connnect to board")
  console.log(err)
  process.exit(1)	
}

let RED
let YELLOW
let GREEN
let SERVO

board.on("ready", () => {
  try {
  	GREEN = new five.Led("GPIO22")
  	YELLOW = new five.Led("GPIO27")
  	RED = new five.Led("GPIO17")
  	SERVO = new five.Servo("GPIO10")

	console.log("getting ready...")
  	RED.blink()
  	YELLOW.blink()
  	GREEN.blink()
  	
 	setTimeout(() => {
	  RED.stop().off()
	  YELLOW.stop().off()
	  GREEN.stop().off()
 	  console.log("program is ready!")
	}, 2000)
  	  
  } catch (err) {
  	console.log("failed to connect to port")
  	console.log(err)
  }

  board.on("exit", () => {
    RED.off()
    GREEN.off()
    YELLOW.off()
  })
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
