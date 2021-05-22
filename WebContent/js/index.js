window.onload = function () {
    displayButton()

    document.getElementById("registerButton").onclick = function () {
        document.getElementById("registerForms").style.display = "block"
    }

    document.getElementById("loginButton").onclick = function () {
        document.getElementById("loginForms").style.display = "block"
    }

    document.getElementById("registerForm").addEventListener('submit', function (event) {
        let fData = new FormData(event.target)
        console.log(fData)
        fetch('/user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(fData)
        })
            .then(function (response) {
                return response.json()
            })
            .then(function (jsonData) {
                window.alert(jsonData["message"])
            })
    })

    document.getElementById("loginForm").addEventListener('submit', function (event) {
        event.preventDefault()

        let fData = new FormData(event.target)
        fetch('/connection', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams(fData)
        })
            .then(function (response) {
                return response.json()
            })
            .then(function (jsonData) {
                window.alert(jsonData["message"])
                if(jsonData["code"] === "200") {
                    document.cookie = "idSession=" + jsonData["idSession"]
                    document.getElementById("loginForms").style.display = "none"
                    displayButton()
                }
            })
    })

    document.getElementById("disconnectButton").onclick = function () {
        fetch('/connection?' + "idSession=" + getCookie("idSession"), {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
                'Accept': 'application/json'
            },
            body: new URLSearchParams()
        })
            .then(function (response) {
                return response.json()
            })
            .then(function (jsonData) {
                window.alert(jsonData["message"])
                if(jsonData["code"] === "200") {
                    console.log(getCookie("idSession"))
                    document.cookie = "idSession=; expires=Thu, 01 Jan 1970 00:00:00 UTC"
                    console.log("id" + getCookie("idSession"))
                }
            })
    }

    document.getElementById("searchButton").onclick = function() {
        let params = new URLSearchParams()
        params.append("search", document.getElementById("search-music").value)
        fetch("/music?" + params.toString())
            .then(res => res.json())
            .then(function (jsonData) {
                document.getElementById("display-music").style.display="none"
                if(jsonData["code"] === "200") {
                    document.getElementById("search-result").style.display="block"
                    const searchResultList = document.getElementById("search-result-list")
                    searchResultList.innerHTML = ""

                    let result = jsonData["result"]["data"]

                    for(let i = 0; i < result.length; i++) {
                        const m = result[i]

                        const form = document.createElement('form')
                        form.name = m["id"]

                        const li = document.createElement('li')

                        const img = document.createElement("img")
                        img.src = m["album"]["cover"]
                        img.width = 50
                        img.height = 50
                        img.style.cursor = "pointer"
                        img.addEventListener("click", function () {
                            displayMusic(form.name)
                        })
                        form.appendChild(img)

                        const infos = document.createElement("ul")

                        const lititle = document.createElement("li")
                        const title = document.createTextNode("Title : " + m["title"] + " ")
                        lititle.appendChild(title)
                        infos.appendChild(lititle)

                        const liartist = document.createElement("li")
                        const artist = document.createTextNode("Artist : " + m["artist"]["name"])
                        liartist.appendChild(artist)
                        infos.appendChild(liartist)

                        const lialbum = document.createElement("li")
                        const album = document.createTextNode("Album : " + m["album"]["title"])
                        lialbum.appendChild(album)
                        infos.appendChild(lialbum)

                        form.appendChild(infos)

                        li.appendChild(form)

                        searchResultList.appendChild(li)
                    }
                }
            })
    }
}

function displayMusic(musicId) {
    console.log(musicId)
    const form = document.createElement('form')
    form.innerHTML = ""
    form.name = musicId
    console.log("step a1")

    let params = new URLSearchParams()
    params.append("music_Id", musicId)
    console.log("step b1")

    fetch("/music?" + params.toString())
        .then(res => res.json())
        .then(function (jsonData) {
            console.log("step c1")
            if(jsonData["code"] === "200") {
                const displayM = document.getElementById("display-music2")
                displayM.innerHTML = ""

                const track = jsonData["result"]
                console.log(track)

                const img = document.createElement("img")
                img.src = track["album"]["cover"]
                img.width = 50
                img.height = 50
                form.appendChild(img)

                const preview = document.createElement("audio")
                preview.controls = true
                preview.volume = 0.2

                const source = document.createElement("source")
                source.src = track["preview"]
                source.type = "audio/mpeg"

                preview.appendChild(source)

                form.appendChild(preview)

                const imgDeezer = document.createElement("img")
                imgDeezer.src = "image/deezerLogo.png"
                imgDeezer.height = 50
                imgDeezer.width = 50
                imgDeezer.onclick = function () {
                    window.open(track["link"])
                }
                imgDeezer.style.cursor = "pointer"
                form.appendChild(imgDeezer)



                const infos = document.createElement("ul")
                const lititle = document.createElement("li")
                const title = document.createTextNode("Title : " + track["title"] + " ")
                lititle.appendChild(title)
                infos.appendChild(lititle)

                const liartist = document.createElement("li")
                const artist = document.createTextNode("Artist : " + track["artist"]["name"])
                liartist.appendChild(artist)
                infos.appendChild(liartist)

                const lialbum = document.createElement("li")
                const album = document.createTextNode("Album : " + track["album"]["title"])
                lialbum.appendChild(album)
                infos.appendChild(lialbum)

                form.appendChild(infos)

                const ul = document.createElement("ul")
                collectMusicLikes(form, ul)
                collectComment(form, ul)

                document.getElementById("search-result").style.display = 'none'

                document.getElementById("display-music").style.display = 'block'

                const form2 = document.createElement('form')
                form2.name = form.name
                const li2 = document.createElement('li')
                console.log("step 3")
                displayM.appendChild(li2)

                displayM.appendChild(form)

                displayM.appendChild(ul)
                console.log("fini")
            }
        })
}

function collectMusicLikes(form, ul) {
    let params = new URLSearchParams()
    params.append("id_Music", form.name)
    fetch("/like_music?" + params.toString())
        .then(res => res.json())
        .then(function (jsonLikes) {
            const li = document.createElement("li")

            let txtNbLikes
            if(jsonLikes["result"].length == 0) {
                txtNbLikes = "Likes : 0"
            } else {
                txtNbLikes = "Likes : " + jsonLikes["result"].length
            }
            const txtNbLikesNode = document.createTextNode(txtNbLikes)
            li.appendChild(txtNbLikesNode)

            const button = document.createElement("button")
            button.innerHTML = "Like Music"
            button.addEventListener("click", function (event) {
                const userId = Promise.resolve(getUserIdFromSession())
                userId.then((id) => {
                    if(id === -1) {
                        window.alert("You need to be connected")
                        return
                    }
                    fetch("/like_music?id_Music=" + form.name + "&id_User=" + id, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json;charset=UTF-8',
                            'Accept': 'application/json'
                        },
                        body: new URLSearchParams()
                    })
                        .then(res => res.json())
                        .then(function(jsonData) {
                            if(jsonData["code"] === "200") {
                                displayMusic(form.name)
                            }
                        })
                })

            })
            li.appendChild(button)

            ul.appendChild(li)
        })
}

function collectComment(form, ul) {
    const userId = Promise.resolve(getUserIdFromSession())

    userId.then((id) => {
        if(id === -1) {
            console.log("Not implemented")
        } else {
            const addComm = document.createElement("input")
            addComm.id = "AddComm"
            addComm.style.width = "50%"
            addComm.defaultValue = "Write a comment"
            form.appendChild(addComm)

            const commBtn = document.createElement("button")
            commBtn.innerHTML = "Post Comment"
            commBtn.onclick = function () {
                let params = new URLSearchParams()
                params.append("id_Music", form.name)
                params.append("id_User", id)
                params.append("msg", addComm.value)
                fetch("/comment?" + params.toString(), {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json;charset=UTF-8',
                        'Accept': 'application/json'
                    }
                })
                    .then(res => res.json())
                    .then(function (jsonData) {
                        window.alert(jsonData["message"])
                        if(jsonData["code"] === "200") {
                            displayMusic(form.name)
                            console.log("added")
                        }
                    })
            }
            form.appendChild(commBtn)
        }
    })

    let params = new URLSearchParams()
    params.append("id_Music", form.name)
    fetch("/comment?" + params.toString())
        .then(res => res.json())
        .then(function (jsonComments) {
            console.log(jsonComments)
            if(jsonComments["result"].length == 0) {
                const li = document.createElement("li")
                const txtComments = document.createTextNode("No comments yet")
                li.appendChild(txtComments)
                ul.appendChild(li)
            } else {
                const result = jsonComments["result"];
                for(let i = 0; i < result.length; i++) {
                    const com = result[i]
                    const login = Promise.resolve(getLoginFromUserId(com["IdUser"]))
                    login.then((log) => {
                        const li = document.createElement("li")
                        const msgForm = document.createElement("form")
                        msgForm.name = com["Id"]

                        const sp = document.createElement("span")
                        sp.style.whiteSpace = "pre-line"

                        if(log === "") {
                            log = "Unkown"
                        }

                        const info = document.createTextNode("Posted by : " + log + " on " + com["Datep"])
                        sp.appendChild(info)
                        msgForm.appendChild(sp)

                        msgForm.append(document.createElement("br"))

                        const sp2 = document.createElement("span")
                        sp2.style.whiteSpace = "pre-line"
                        const msg = document.createTextNode(com["Msg"])
                        sp2.appendChild(msg)
                        msgForm.appendChild(sp2)

                        msgForm.append(document.createElement("br"))

                        msgForm.appendChild(document.createTextNode("Likes : " + com["Likes"]))

                        const likeButton = document.createElement("img")
                        likeButton.src = "image/likeButton.png"
                        likeButton.width = 12
                        likeButton.height = 12
                        likeButton.style.marginLeft = "10px"
                        likeButton.style.cursor = "pointer"
                        likeButton.onclick = function() {
                            const userId = Promise.resolve(getUserIdFromSession())
                            userId.then((id) => {
                                if (id == -1) {
                                    window.alert("You need to be connected")
                                }
                                let params = new URLSearchParams()
                                params.append("id_Comment", msgForm.name)
                                params.append("id_User", id)
                                fetch("/like_comment?" + params.toString(), {
                                    method: 'POST',
                                    headers: {
                                        'Content-Type': 'application/json;charset=UTF-8',
                                        'Accept': 'application/json'
                                    }
                                })
                                    .then(res => res.json())
                                    .then(function (jsonData) {
                                        if (jsonData["code"] === "200") {
                                            displayMusic(form.name)
                                        }
                                    })
                            })
                        }
                        msgForm.appendChild(likeButton)


                        li.appendChild(msgForm)

                        ul.appendChild(li)
                    })


                }

            }
        })
}

function displayButton() {
    const idSession = getCookie("idSession")
    document.getElementById("registerButton").style.display = "block"
    document.getElementById("loginButton").style.display = "block"
    document.getElementById("disconnectButton").style.display = "none"
    if(idSession != "") {
        fetch("/connection?idSession=" + idSession)
            .then(res => res.json())
            .then(function (jsonData) {
                if(jsonData["code"] === "200") {
                    document.getElementById("registerButton").style.display = "none"
                    document.getElementById("loginButton").style.display = "none"
                    document.getElementById("disconnectButton").style.display = "block"
                }
            })
    }
}

function getUserIdFromSession() {
    const idSession = getCookie("idSession")
    let params = new URLSearchParams()
    params.append("idSession", idSession)
    if (idSession != "") {
        return fetch("/connection?" + params.toString())
            .then(res => res.json())
            .then(function (jsonData) {
                if (jsonData["code"] === "200") {
                    return jsonData["userId"]
                } else {
                    return -1
                }
            })
    }
    return -1
}

function getLoginFromUserId(id) {
    let params = new URLSearchParams()
    params.append("id", id)
    return fetch("/user?" + params.toString())
        .then(res => res.json())
        .then(function (jsonData) {
            if(jsonData["code"] === "200") {
                return jsonData["login"]
            } else {
                return ""
            }
        })
    return ""
}

function getCookie(cname) {
    const name = cname + "=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const ca = decodedCookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}