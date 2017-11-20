$(function(){
	
	//回头部
	$('.back_to_top').hide();
	$(window).scroll(function(){
		if($(window).scrollTop()>450){
			$('.back_to_top').fadeIn();	
		}else{
			$('.back_to_top').fadeOut();	
		}
	});	
	
	$('.back_to_top').click(function(){
		$('html,body').animate({scrollTop:0},600);
		return false;
	});
	
	//导航下拉
	$('.nav li').hover(function(){
		$(this).stop().addClass('hover2').find('.two_nav').show();	
	},function(){
		$(this).stop().removeClass('hover2').find('.two_nav').hide();		
	});
	
	$('.lang').hover(function(){
		$(this).find('ul').slideDown();	
	},function(){
		$(this).find('ul').slideUp();	
	});
	
	//自适应
	if($(window).width()<1800){
		$('.head').addClass('head_x');	
	}else{
		$('.head').removeClass('head_x');		
	}
	
	if($(window).width()>1656){
		$('.w_bfb').addClass('w_i');	
	}else{
		$('.w_bfb').removeClass('w_i');		
	}
	
	$(window).resize(function(){
		
		
		if($(window).width()<1800){
			$('.head').addClass('head_x');	
		}else{
			$('.head').removeClass('head_x');		
		};
		
		
		if($(window).width()>1656){
		$('.w_bfb').addClass('w_i');	
		}else{
			$('.w_bfb').removeClass('w_i');		
		};
	});
	
	
	$('.pro_i li').hover(function(){
		$(this).find('.pro_nr').animate({bottom:'0'},300);	
	},function(){
		$(this).find('.pro_nr').animate({bottom:'-100%'},300);		
	});
	
	//新闻列表
	$('.news_list dl').hover(function(){
		$(this).find('.bj').fadeIn();
		$(this).find('.share_lb').slideDown();
	},function(){
		$(this).find('.share_lb').slideUp();
		$(this).find('.bj').fadeOut();
	});
	
	//底部二维码
	$('.erm_li').hover(function(){
		$(this).find('.erm').slideDown();	
	},function(){
		$(this).find('.erm').slideUp();	
	});
	
})

//选项卡
function setTab(name,cursel,n){
for(i=1;i<=n;i++){
var menu=document.getElementById(name+i);
var con=document.getElementById("con_"+name+"_"+i);
menu.className=i==cursel?"hover":"";
con.style.display=i==cursel?"block":"none";
}
}


function navFixed(){
	var oConbox = $('.main_ny .con_box');
	var oNynay = $('.main_ny .ny_nav');
	if (oConbox.size() > 0 && oNynay.height() < $(window).height() - oConbox.offset().top) {
		oNynay.css('position','fixed');
		var oConboxTop = oConbox.offset().top;
		var maxTop = oConbox.height() - oNynay.height();
		$(window).on('scroll',function(){
			var sTop = $(window).scrollTop();
			console.log(maxTop,sTop);
			if (sTop > maxTop){
				oNynay.css('top',maxTop - sTop + oConboxTop);
			}else{
				oNynay.css('top',oConboxTop);
			}
		})
	}
}

// window.onload = function(){
// 	navFixed();
// }