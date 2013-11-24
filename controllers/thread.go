package controllers

import (
    "github.com/gorilla/mux"
    "database/sql"
    "time"
    "net/http"
    "strconv"
    "github.com/stevenleeg/gobb/utils"
    "github.com/stevenleeg/gobb/models"
)

func Thread(w http.ResponseWriter, r *http.Request) {
    board_id_str := mux.Vars(r)["board_id"]
    board_id, _ := strconv.Atoi(board_id_str)
    err, board := models.GetBoard(board_id)

    post_id_str := mux.Vars(r)["post_id"]
    post_id, _ := strconv.Atoi(post_id_str)
    err, op, posts := models.GetThread(post_id)

    if r.Method == "POST" {
        db := models.GetDbSession()
        title := r.FormValue("title")
        content := r.FormValue("content")

        current_user := utils.GetCurrentUser(r)
        if current_user == nil {
            http.NotFound(w, r)
            return
        }

        post := models.NewPost(current_user, board, title, content)
        post.ParentId = sql.NullInt64{ int64(post_id), true }
        op.LatestReply = time.Now()
        db.Insert(post)
        db.Update(op)

        err, op, posts = models.GetThread(post_id)
    }

    if err != nil {
        http.NotFound(w, r)
        return
    }

    utils.RenderTemplate(w, r, "thread.html", map[string]interface{} {
        "board": board,
        "op": op,
        "posts": posts,
    })
}
