// Copyright 2021 hamzashezad. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"net/http"
	"time"

	"miniflux.app/http/request"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/model"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

func (h *handler) showTodayPage(w http.ResponseWriter, r *http.Request) {
	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)

	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	now := time.Now()
	y, m, d := now.Date()

	today := time.Date(y, m, d, 0, 0, 0, 0, now.Location());

	offset := request.QueryIntParam(r, "offset", 0)
	builder := h.store.NewEntryQueryBuilder(user.ID)
	builder.AfterDate(today)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithGloballyVisible()
	countTodayUnread, err := builder.CountEntries()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	if offset >= countTodayUnread {
		offset = 0
	}

	builder = h.store.NewEntryQueryBuilder(user.ID)
	builder.AfterDate(today)
	builder.WithStatus(model.EntryStatusUnread)
	builder.WithOrder(model.DefaultSortingOrder)
	builder.WithDirection(user.EntryDirection)
	builder.WithOffset(offset)
	builder.WithLimit(user.EntriesPerPage)
	builder.WithGloballyVisible()
	entries, err := builder.GetEntries()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	view.Set("entries", entries)
	view.Set("pagination", getPagination(route.Path(h.router, "today"), countTodayUnread, offset, user.EntriesPerPage))
	view.Set("menu", "today")
	view.Set("user", user)
	view.Set("countTodayUnread", countTodayUnread)
	view.Set("countUnread", h.store.CountUnreadEntries(user.ID))
	view.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(user.ID))
	view.Set("hasSaveEntry", h.store.HasSaveEntry(user.ID))

	render := view.Render("today_entries")

	html.OK(w, r, render)
}
