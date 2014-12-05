$(function(){
    $('.status').click(function(){
        var code_id = $(this).siblings().first().children().first().html().slice(1);
        console.log(code_id);
        var data = {};
		$('.loader').removeClass('disable');
		$('.loader').addClass('active');
        data.id = code_id;
        $.getJSON("/Code/GetPanic",data,function(data,textStatus,jqXHR){
            if(data.status){
                var panic = data.panic;
                $('#check').empty(); 
				$('#check').html("<p>"+panic+"</p>");
			    $('.loader').addClass('disable');
			    $('.loader').removeClass('active');
			    $('.ui.modal').modal('show');
            }else{
			    $('.loader').addClass('disable');
			    $('.loader').removeClass('active');
                console.log(data.error);
            };
        });
    });
	$('.view').click(function(){
		var code_id = $(this).siblings().first().children().first().html().slice(1);
		console.log(code_id);
		var data = {};
		data.id = code_id;
		$('.loader').removeClass('disable');
		$('.loader').addClass('active');
		$.getJSON("/Code/View",data,function(data,textStatus,jqXHR){
			if(data.status){
				var view = data.code;
				console.log(view);
				var pre = "<pre class=\"prettyprint\"></pre>";
				$('#check').empty();
				$('#check').siblings().eq(1).text("Code");
				$(pre).appendTo('#check');
				$(".prettyprint").text(view);
				prettyPrint();
				$('.loader').addClass('disable');
				$('.loader').removeClass('active');
				$('.ui.modal').modal('show');
			}else{
				$('.loader').addClass('disable');
				$('.loader').removeClass('active');
				$('#check').html("<p>"+data.error+"</p>")
				console.log(data.error);
			}

		});		
	});
	$('.id').click(function(){
		var code_id = $(this).children().first().html().slice(1);
		console.log("code is "+code_id);
		var data ={};
		data.id = code_id;
		$('.loader').removeClass('disable');
		$('.loader').addClass('active');
		$.getJSON("/Code/Check",data,function(data,textStatus,jqXHR){
			if(data.status){
				console.log(data.report);
				tests = data.report.Tests
				var contents = [];
				var divider = "<div class=\"ui section divider\"></div>"
				for(var i=0;i<tests.length;i++){
					var nth = i;
					var input = tests[i].In.replace(/\n/g,"<br/>");
					input = "Input:<br/>" + input;
					var output = tests[i].Out.replace(/\n/g,"<br/>");
					output = "Output:<br/>" + output;	
					var t = "<p>#"+i+"<br/>"+input + output+"</p>";
					contents.push(t);
					console.log($('#check').html());
				}
				var content = contents.join(divider);
				$('#check').empty();
				$('#check').siblings().eq(1).text("Tests");
				$('#check').html(content);
				$('.loader').addClass('disable');
				$('.loader').removeClass('active');
				$('.ui.modal').modal('show');
			}else{
				$('.loader').addClass('disable');
				$('.loader').removeClass('active');
				$('#check').html("<p>"+data.error+"</p>")
				console.log(data.error);
			}
		});
	});
});
