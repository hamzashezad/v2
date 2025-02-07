// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/logger"
	"miniflux.app/model"
	"miniflux.app/ui/form"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
	"miniflux.app/validator"
)

func (h *handler) saveCategory(w http.ResponseWriter, r *http.Request) {
	loggedUser, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	categoryForm := form.NewCategoryForm(r)

	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)
	view.Set("form", categoryForm)
	view.Set("menu", "categories")
	view.Set("user", loggedUser)
	view.Set("countTodayUnread", h.store.CountTodayUnreadEntries(loggedUser.ID))
	view.Set("countUnread", h.store.CountUnreadEntries(loggedUser.ID))
	view.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(loggedUser.ID))

	categoryRequest := &model.CategoryRequest{Title: categoryForm.Title}

	if validationErr := validator.ValidateCategoryCreation(h.store, loggedUser.ID, categoryRequest); validationErr != nil {
		view.Set("errorMessage", validationErr.TranslationKey)
		html.OK(w, r, view.Render("create_category"))
		return
	}

	if _, err = h.store.CreateCategory(loggedUser.ID, categoryRequest); err != nil {
		logger.Error("[UI:SaveCategory] %v", err)
		view.Set("errorMessage", "error.unable_to_create_category")
		html.OK(w, r, view.Render("create_category"))
		return
	}

	html.Redirect(w, r, route.Path(h.router, "categories"))
}
