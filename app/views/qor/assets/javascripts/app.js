"use strict";$(function(){$(document).on("click.qor.alert",'[data-dismiss="alert"]',function(){$(this).closest(".qor-alert").removeClass("qor-alert__active")}),setTimeout(function(){$('.qor-alert[data-dismissible="true"]').removeClass("qor-alert__active")},5e3)}),$(function(){$(document).on("click",".qor-dialog--global-search",function(e){e.stopPropagation(),$(e.target).parents(".qor-dialog-content").length||$(e.target).is(".qor-dialog-content")||$(".qor-dialog--global-search").remove()}),$(document).on("click",".qor-global-search--show",function(e){e.preventDefault();var a=$(this).data(),o=window.Mustache.render('<div class="qor-dialog qor-dialog--global-search" tabindex="-1" role="dialog" aria-hidden="true"><div class="qor-dialog-content"><form action=[[actionUrl]]><div class="mdl-textfield mdl-js-textfield" id="global-search-textfield"><input class="mdl-textfield__input ignore-dirtyform" name="keyword" id="globalSearch" value="" type="text" placeholder="" /><label class="mdl-textfield__label" for="globalSearch">[[placeholder]]</label></div></form></div></div>',a);$("body").append(o),window.componentHandler.upgradeElement(document.getElementById("global-search-textfield")),$("#globalSearch").focus()})}),$(function(){var l=[],s="qoradmin_menu_status",e=localStorage.getItem(s);e&&e.length&&(l=e.split(",")),$(".qor-menu-container").on("click","> ul > li > a",function(){var e=$(this),a=e.parent(),o=e.next("ul"),t=a.attr("qor-icon-name");o.length&&(o.hasClass("in")?(l.push(t),a.removeClass("is-expanded"),o.one("transitionend",function(){o.removeClass("collapsing in")}).addClass("collapsing").height(0)):(l=_.without(l,t),a.addClass("is-expanded"),o.one("transitionend",function(){o.removeClass("collapsing")}).addClass("collapsing in").height(o.prop("scrollHeight"))),localStorage.setItem(s,l))}).find("> ul > li > a").each(function(){var e=$(this),a=e.parent(),o=e.next("ul"),t=a.attr("qor-icon-name");o.length&&(o.addClass("collapse"),a.addClass("has-menu"),-1!=l.indexOf(t)?o.height(0):(a.addClass("is-expanded"),o.addClass("in").height(o.prop("scrollHeight"))))});var a=$(".qor-page > .qor-page__header"),o=$(".qor-page > .qor-page__body"),t=a.find(".qor-page-subnav__header").length?96:48;a.length&&(a.height()>t&&o.css("padding-top",a.height()),$(".qor-page").addClass("has-header"),$("header.mdl-layout__header").addClass("has-action"))}),$(function(){$(".qor-mobile--show-actions").on("click",function(){$(".qor-page__header").toggleClass("actions-show")})}),$(function(){var p,m,q=$("body"),f="is-selected",b=function(){return q.hasClass("qor-bottomsheets-open")};function _(e){$("[data-url]").removeClass(f),e&&e.length&&e.addClass(f)}q.qorBottomSheets(),q.qorSlideout(),p=q.data("qor.slideout"),m=q.data("qor.bottomsheets"),$(document).on("click.qor.openUrl","[data-url]",function(e){var a,o=$(this),t=$(e.target),l=o.hasClass("qor-button--new"),s=o.hasClass("qor-button--edit"),n=(o.is(".qor-table tr[data-url]")||o.closest(".qor-js-table").length)&&!o.closest(".qor-slideout").length,r=o.data(),i=r.openType,d=o.parents(".qor-theme-slideout").length,c=o.closest(".qor-slideout").length,h=o.hasClass("qor-action-button")||o.hasClass("qor-action--button");if(e.stopPropagation(),!(o.data("ajax-form")||t.closest(".qor-table--bulking").length||t.closest(".qor-button--actions").length||!t.data("url")&&t.is("a")||n&&b()))if("window"!=i){var u,g;if("new_window"!=i)return h&&(u=$(".qor-js-table tbody").find(".mdl-checkbox__input:checked"),g=[],(a=!!u.length&&(u.each(function(){g.push($(this).closest("tr").data("primary-key"))}),g))&&(r=$.extend({},r,{actionData:a}))),r.$target=t,r.method&&"GET"!=r.method.toUpperCase()?void 0:("bottomsheet"!=i&&!h||"slideout"==i?"slideout"==i||n||l&&!b()||s?"slideout"==i||d?o.hasClass(f)?(p.hide(),_()):(p.open(r),_(o)):window.location.href=r.url:q.hasClass("qor-slideout-open")||l&&b()?m.open(r):d?p.open(r):m.open(r):h&&!a&&o.closest('[data-toggle="qor.action.bulk"]').length&&!c?window.QOR.qorConfirm(r.errorNoItem):m.open(r),!1);window.open(r.url,"_blank")}else window.location.href=r.url})}),$(function(){var l=window.location;$(".qor-search").each(function(){var e=$(this),a=e.find(".qor-search__input"),o=e.find(".qor-search__clear"),t=!!a.val();e.closest(".qor-page__header").addClass("has-search"),$("header.mdl-layout__header").addClass("has-search"),o.on("click",function(){a.val()||t?"?"==l.search.replace(new RegExp(a.attr("name")+"\\=?\\w*"),"")?l.href=l.href.split("?")[0]:l.search=l.search.replace(new RegExp(a.attr("name")+"\\=?\\w*"),""):e.removeClass("is-dirty")})})});
$.ajaxSetup({
    statusCode: {
        406: function () {
            $.ajax({
                url: 'http://localhost:8080/refresh',
                type: 'POST',
                beforeSend: function(xhr){
                    xhr.withCredentials = true;
                },
                success: (data) => {
                    console.log(data)
                }
            });
        }
    }
});
$( document ).ajaxError(function( event, jqxhr, settings, thrownError ) {
    if ( jqxhr.status === 406 ) {
        $.ajax({
            url: 'http://localhost:8080/refresh',
            type: 'POST',
            beforeSend: function(xhr){
                xhr.withCredentials = true;
            },
            success: (data) => {
                console.log(data)
            }
        });
    }
  });