function StringToJson(data: string): Array<object> {

  let result = [];
  if (data.length > 2) {
    if (!data.startsWith("[")) {
      data = "[" + data + "]";
    }
    result = JSON.parse(data);
  }
  return result;
}

export {
  StringToJson
}