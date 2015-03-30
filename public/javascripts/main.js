$(document).on("click", "a.sign-out", function(e) {
  e.preventDefault()
  $link = $(this)

  $.ajax({
    method: "DELETE",
    url: $link.attr("href"),
    dataType: "json",
    success: function(json) {
      window.location.replace(json.location);
    }
  })
})
