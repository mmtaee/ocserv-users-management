// const Rules = {
//     required: (v : any) => !!v || "required",
// };

// export default Rules;

function required(v: any): any {
  return !!v || "Required";
}

function number(v: any): any {
  return v && isNaN(v) ? "Number Required" : true;
}

function ip(v: any): any {
  var ipformat =
    /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  return v && !v.match(ipformat) ? "IP Format Required" : true;
}

function ipOrRange(v: any): any {
  var ipformat =
    /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?:\/(?:3[0-2]|[1-2]?[0-9]))?$|^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.)?(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?:\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){2}$/;
  return v && !v.match(ipformat) ? "IP Format Required" : true;
}

export { required, number, ip , ipOrRange};
