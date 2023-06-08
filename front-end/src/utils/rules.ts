// const Rules = {
//     required: (v : any) => !!v || "required",
// };

// export default Rules;

function required(v: any): any {
    return !!v || "required"
}

export {
    required,
}