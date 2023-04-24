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
            linkBtnProxy(offcanvasBody)
            avatarUploadProxy(offcanvasBody)
            updateProfixProxy(offcanvasBody)


            // 判断点击了哪一个元素
            if (event.target.dataset.tag === "myurls") {
                window.offcanvas.title.innerHTML = "MyURLs"
                // 请求
                fetchList()
            } else {
                window.offcanvas.title.innerHTML = "Update Profile"
                // 请求
                fetch("/profile", {
                    method: "get",
                }).then(response => response.json()).then(data => {
                    errorHandler(data)
                    if (data.status == 200) {
                        let user = data.data
                        offcanvasBody.innerHTML = `
                        <div class="profile">
                            <form action="#" method="post">
                                <div class="user_name">
                                    <input type="text" name="name" value="${user.name}" />
                                </div>
                                <div class="user_avatar">
                                    <input type="file" class="iconfont" name="file"/>  
                                    <div class="box">
                                        <img id="user_avatar" src="${user.avatar}" />
                                    </div>
                                </div> 
                                <div class="submit">
                                    <input type="submit" value="Save">
                                </div>
                            </form>
                        </div>

                        <div class="user_menu">
                            <ul>
                                <li class="logout">
                                    <a href="/logout">退出登录</a>
                                </li>
                            </ul>
                        </div>
                    `
                    }
                })
            }
        })
    })
})

function errorHandler(data) {
    if (data.status != 200) {
        window.createToast("error", data.error)
        let timer = setTimeout(() => {
            if (300 <= data.status < 400 && data.redirect) {
                location.href = data.redirect
                return
            }
            clearTimeout(timer)
        }, 3000)
        return
    }
}

function linkBtnProxy(offcanvasBody) {
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
        errorHandler(data)
        let links = data.data
        let html = ""
        if (links.length === 0) {
            return
        }
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
        errorHandler(data)
    })
}

function avatarUploadProxy(offcanvasBody) {
    offcanvasBody.addEventListener("change", function(event) {
        const file = event.target.files[0]
        let url = URL.createObjectURL(file)
        document.getElementById("user_avatar").src = url
    })
}

function updateProfixProxy(offcanvasBody) {
    offcanvasBody.addEventListener("submit", function(event) {
        event.preventDefault()
        event.stopPropagation()

        let formData = new FormData(event.target)
        let tokenInput = document.getElementById("fetch_csrf_token")
        formData.append("csrf_token", tokenInput.value)

        fetch("/profile", {
            method: "post",
            body: formData,
        }).then(response => response.json()).then(data => {
            errorHandler(data)
            if (data.status == 200) {
                window.createToast("success", "更新成功")
            }
        })
    })
}