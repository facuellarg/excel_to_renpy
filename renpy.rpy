define John = Character("John")
define Tom = Character("Tom")

label start:

  show John happy at left
  John "Hello"
  show Tom happy at left
  Tom "How are you?"
  menu:
    "option1":
      jump otherLabel
    "option2"
    "option3"
  hide John

  scene imageScene
  show Tom happy at left
  Tom "Hello after scene"
  show Tom angry at left
  Tom "I am angry now"

label otherLabel:

  Tom "Hello from another label"

