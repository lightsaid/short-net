// 右边侧栏功能模块
window.addEventListener("load", function () {
    fetchList.pageInfo = {
        pageIndex: 0,
        pageSize: 10
    }
    window.offcanvas.openBtns.forEach((btn) => {
        btn.addEventListener("click", function (event) {
            // console.log(event.target.dataset.tag)

            let offcanvasBody = document.querySelector(".offcanvas-body")

            // 代理事件
            linkBtnPorxy(offcanvasBody)

            // 判断点击了哪一个元素
            if (event.target.dataset.tag === "myurls") {
                window.offcanvas.title.innerHTML = "MyURLs"
                // 请求
                fetchList()
            } else {
                window.offcanvas.title.innerHTML = "Profile"
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
})

function linkBtnPorxy(offcanvasBody) {
    offcanvasBody.addEventListener("click", function(event) {
        console.log(event.target.dataset.btn)

        let btn = event.target.dataset.btn
        if (btn == "copy") {
            let target = event.target
            let input = target.querySelector("input")
            input.select()
            window.document.execCommand("copy")
            window.createToast("success", "已复制到剪贴板")
        }

        if (btn=="delete") {
            fetchDelete(event.target.dataset.id)
        }

        if (btn=="loading") {
            fetchList.pageInfo.pageIndex++
            fetchList()
        }
    })
}


function fetchList() {

    let offcanvasBody = document.querySelector(".offcanvas-body")

    const { pageIndex, pageSize } = fetchList.pageInfo

    let formData = new FormData()
    formData.append("page", pageIndex)
    formData.append("size", pageSize)

    // NOTE: 获取不到
    // formData.append("csrf_token", "{{.CSRFToken}}")
    let tokenInput = document.getElementById("fetch_csrf_token")
    formData.append("csrf_token", tokenInput.value)

    return fetch("/short/list", {
        method: "post",
        body: formData,
    }).then(response => response.json()).then(data => {
        let links = data.data
        let html = ""
        links.forEach((link)=>{
            let short_url = `${location.origin}/${link.short_hash}`
            let expired_at =  new Date(link.expired_at).getTime()
            let current_at = Date.now()
            let expire_text = "Expire"
            if (current_at < expired_at) {
                expire_text = "Valid"
            }
            html += `
                <ul class="links">
                    <li class="short_url">
                        <a target="_blank" href="${short_url}">${short_url}</a>
                    </li>
                    <li class="long_url">
                        <a target="_blank" href="${link.long_url}">${link.long_url}</a>
                    </li>
                    <li class="info">
                        <div>
                            <span>${link.click} click</span>
                            <span style="font-size: 12px;"> &nbsp; | &nbsp;</span>
                            <span>${link.updated_at}</span>
                        </div>
                    </li>
                    <li class="opt">
                        <div>
                            <button data-btn="expire">${expire_text}</button>
                            <button class="copy" data-btn="copy">Copy<input style="position: fixed; top: -1000000px;" type="text" value="${short_url}"></button>
                            <button class="delete" data-btn="delete" data-id="${link.id}">Delete</button>
                        </div>
                    </li>
                </ul>
            `
        })
        if (links.length < 10){
            html += ` <div class="loading"> <span  disable data-btn="loading">no more</span></div>`
        }else{
            html += `<div class="loading"> <button  data-btn="loading">Loading</button></div>`
        }
        offcanvasBody.innerHTML = html


    }).catch(err => {
        let msg = "服务内部错误"
        // 信息有异常
        if (err.message && err.message.length > 20) {
            window.createToast("error", msg)
        } else {
            window.createToast("error", err.message || msg)
        }
        console.error("createLink: ", err)
    })
}


function fetchDelete(linkId) {
    let formData = new FormData()
    let tokenInput = document.getElementById("fetch_csrf_token")
    formData.append("csrf_token", tokenInput.value)
    formData.append("link_id", linkId)

    return fetch("/short/delete", {
        method: "post",
        body: formData,
    }).then(response => response.json()).then(data => {
        if (data.status == 200) {
            window.createToast("success", "删除成功")
            fetchList()
        }
    }).catch(err => {
        let msg = "服务内部错误"
        // 信息有异常
        if (err.message && err.message.length > 20) {
            window.createToast("error", msg)
        } else {
            window.createToast("error", err.message || msg)
        }
        console.error("createLink: ", err)
    })
}