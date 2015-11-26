package main

func (operation *Operation_Init) Init_Default_Run(flags []string) (bool, map[string]string) {
  key := "bare"
  if len(flags)>0 {
    key = flags[0]
    flags = flags[1:]
  }

  switch key {
    case "bare":
      return true, operation.Init_Default_Bare()
    case "starter":
      return true, operation.Init_Default_Starter()
  }

  return false, map[string]string{}
}

func (operation *Operation_Init) Init_Default_Bare() map[string]string {
  return map[string]string{
    
  }
}

func (operation *Operation_Init) Init_Default_Starter() map[string]string {
  return map[string]string{

  }
}
