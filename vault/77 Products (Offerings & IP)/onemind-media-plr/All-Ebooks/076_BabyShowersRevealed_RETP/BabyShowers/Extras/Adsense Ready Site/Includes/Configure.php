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
	
	$Category				= "Baby Shower";
	$NewsFeed				= "Baby+Shower";
			
	$Keywords				= "Baby Shower, Baby Shower Cake, Baby Shower Cake Ideas, Baby Shower Cakes, Baby Shower Clip Art, Baby Shower Decorations, Baby Shower Etiquette, Baby Shower Favors, Baby Shower Food, Baby Shower Footprint Invitations, Baby Shower Game Ideas, Baby Shower Games, Baby Shower Gift, Baby Shower Gift Basket, Baby Shower Gift Idea, Baby Shower Gift Ideas,Baby Shower Gifts,Baby Shower Graphics, Baby Shower Idea, Baby Shower Ideas, Baby Shower Invitation, Baby Shower Invitation Ideas, Baby Shower Invitation Wording, Baby Shower Invitations, Baby Shower Invites, Baby Shower Party Favors, Baby Shower Planning, Baby Shower Poems, Baby Shower Recipes, Baby Shower Supplies, Baby Shower Theme, Baby Shower Themes, Discounted Baby Shower Supplies, Free Baby Shower Games, Free Baby Shower Invitation, Free Baby Shower Invitations, Free Printable Baby Shower Games, Free Printable Baby Shower Invitations, Homemade Baby Shower Favors, Inexpensive Baby Shower Gifts, Planning A Baby Shower, Unique Baby Shower Gift Ideas, Unique Baby Shower Gifts, Unique Baby Shower Invitations, Wording For Baby Shower Invitations";
	
	$GooglePubID			= ""; 
	$GoogleChannel			= ""; 
	$GoogleLinkColor		= "";
	$GoogleURLColor			= "FF93BA";	

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

?>