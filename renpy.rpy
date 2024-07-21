define John = Character("John")

label start:

  John "Hello"
  John "How are you?"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"

  scene imageScene
