// 右边侧栏功能模块
window.addEventListener("load", function () {
    window.onload = function () {
        window.offcanvas.openBtns.forEach((btn) => {
            btn.addEventListener("click", function (event) {
                // console.log(event.target.dataset.tag)
    
                let offcanvasBody = document.querySelector(".offcanvas-body")
    
                // 判断点击了哪一个元素
                if (event.target.dataset.tag === "myurls") {
                    console.log(offcanvasBody)
                    offcanvasBody.innerHTML = `
                    <div>
                        <p>锻炼</p>    
                        <p>锻炼</p>    
                        <p>锻炼</p>    
                        <p>锻炼</p>    
                    </div>
                `
                } else {
                    offcanvasBody.innerHTML = `
                    <div>
                        <p>锻炼2</p>    
                        <p>锻炼2</p>    
                        <p>锻炼2</p>    
                        <p>锻炼2</p>    
                    </div>
                `
                }
            })
        })
    }
})
