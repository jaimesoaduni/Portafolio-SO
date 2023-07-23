const others = document.getElementsByClassName("other")

{
    let i = 1
    for (let elem of others) {
        if (i % 2 == 1) {
            console.log(elem)
            elem.style.background = "#c7ccd4"
        }
        i++
    }

}