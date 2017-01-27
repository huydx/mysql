list = make([]item, 0)

for {
  token = nextToken()
  if isType1(token) {
    list = append(list, item1(token))
  } else isType2(token) {
    list = append(list, item2(token))
  } .....
}

return list
