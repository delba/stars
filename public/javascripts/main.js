$(document).on("click", "a.sign-out", function(e) {
  e.preventDefault();
  $link = $(this);

  $.ajax({
    method: "DELETE",
    url: $link.attr("href"),
    dataType: "json",
    success: function(json) {
      window.location.replace(json.location);
    }
  })
})

$(document).on("click", "a.star", function(e) {
  e.preventDefault();
  $link = $(this);

  $.ajax({
    method: "PUT",
    url: $link.attr("href"),
    dataType: "json",
    success: function(json) {
      // TODO change star class to unstar
      console.log(json);
    }
  })
})

$(document).on("click", "a.unstar", function(e) {
  e.preventDefault();
  $link = $(this);

  $.ajax({
    method: "DELETE",
    url: $link.attr("href"),
    dataType: "json",
    success: function(json) {
      // TODO change unstar class to star
      console.log(json)
    }
  })
})
