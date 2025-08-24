const arr = [1,12,2,23,4,4,5]
const arr1 = [1,12,2,23,4,4,5]

const arr2 = [...arr, arr1]

const [a,b, ...rest] = arr 

// console.log(arr2)
console.log("a" , a , "rest", rest )