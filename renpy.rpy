define John = Character("John")
define Tom = Character("Tom")

label start:

  John "Hello"
  Tom "How are you?"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"

  scene imageScene

label otherLabel:

  Tom "Hello from another label"

