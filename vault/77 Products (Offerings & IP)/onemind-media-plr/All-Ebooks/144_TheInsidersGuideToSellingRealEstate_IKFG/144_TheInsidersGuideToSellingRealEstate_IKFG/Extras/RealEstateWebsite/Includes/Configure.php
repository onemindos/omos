<?php 

 	header('Expires: Mon, 26 Jul 1997 05:00:00 GMT');
	header('Cache-Control: no-store, no-cache, must-revalidate');
	header('Cache-Control: post-check=0, pre-check=0', FALSE);
	header('Pragma: no-cache');

	$UserName				= "admin";
	$Password				= "password";
	
	$Domain					= "http://".$_SERVER["HTTP_HOST"].rtrim(str_replace($DOCUMENT_ROOT, "", dirname(dirname($PHP_SELF))), "/")."/";
	$SiteName				= "";
		
	$MainTitle				= "";
	$SubTitle				= "";

	$Category				= "Real Estate";
	$NewsFeed				= "Real+Estate";
	
	$Keywords				= "Advantages Of Buying A Property In France, Best Way To Sell Property, Buying A Home After Bankruptcy, Buying A Home In Houston Texas, Buying A Home In Las Vegas, Buying A Home In Memphis Tennessee, Buying A Home In Naples , Buying A Home In San Antonio, Buying A Home With Bad Credit , Buying A Tenerife Property, Buying Homes In Henderson, Buying Homes In Long Island, Buying Homes In Portugal , Buying Homes With Resale Value, Buying Properties In Bulgaria, Buying Properties In Spain, Buying Property Abroad, Buying Property In Abbotsford Canada, Buying Property In Cyprus, Buying Property In Dubai, Buying Property In South France, Buying Property In Spain, Buying Property In Turkey, Buying You Dream House In Denver, Buying Your First Home In Las Vegas, Essential Tips On Buying Property, Guide To Buying Property In France, Home Buying Surnames, Homes Of Your Dreams, How To Buy Property In Spain, How To Sell A Home In Conway, Arkansas, Information On Buying Property In SPain , Selling A Home In Colorado, Selling A Home In FLorida, Selling Home By Owner, Selling Homes In Chicago Illinois, Selling Homes In Sarasota Florida, Selling Homes In Summerlin, Selling Homes In Victoria , Selling Property In Long Island, Selling You Home In Summerlin, Selling Your Bonita Spring Homes, Tidbits On Buying And Selling Homes, Tips For Buying A Bulgaria Property, Tips On Buying A Home, Tips On Buying A Homesite, Tips On Buying And Selling A Home, Tips On Buying Property In SPain, Tips On Selling Homes, Tips On Selling Homes In Las Vegas, Tips To Remember When Selling Homes";
	
	$GooglePubID			= "";
	$GoogleChannel			= "";
	$GoogleLinkColor		= "A40F15";
	$GoogleURLColor			= "D76062";

	$AmazonID				= "";
	$AmazonMode				= "Books";
	$AmazonSearch			= "Real Estate";
	$AmazonLinkColor		= "A40F15";
	$AmazonTextColor		= "D76062";
	$DisplayAmazon			= "";

	$CopyRight				= "All images, text and code found on this web site are &copy; ".date("Y").",".$SiteName." - All rights reserved worldwide...";

	$Name					= "";
	$Address1				= "";
	$Address2				= "";
	$City					= "";
	$State					= "";
	$ZipCode				= "";
	$Country				= "";
	$Telephone				= "";
	$EmailAddress			= "";

	$ShowGoogle				= "checked";
	$ShowFireFox			= "checked";
	$ShowSearch				= "checked";

	$RotateArticles			= "checked";	

?>