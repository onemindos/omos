<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | Ebay Training Page 1</title>
		<META http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Creator" content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright" content="<?=$CopyRight;?>">
		<meta name="Description" content="Ebay Training Related Links - Page 1 - ">
		<meta name="Keywords" content="<?=$Keywords;?>">
		<meta name="Distribution" content="global">
		<meta name="Publisher" content="<?=$Domain;?>">
		<meta name="Rating" content="General">
		<meta name="Revisit-after" content="5 days">
		<meta name="Robots" content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0">
			<tr>
				<td class="tdHeader" colspan="3">
					<p style="margin-bottom: 5; margin-left: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_728x15_TextLinks.php"); ?>
					<hr>
				</td>
			</tr>
			<tr>
				<td class="tdleft">
					<br>
					<?php include("../Includes/Navigation.php"); ?>
				</td>
				<td class="tdContent">
					<br>
					<p>
						<b>Ebay Training Related Links - Page 1 </b>
						<br>
						<b><small> </small></b>
						<br>
						<br>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<ul class='ulCategory'><li><A href="http://www.jackwalters.com/" target='_blank'>eBay training and Business Consulting by Jack Walters Certified eBay ... </A><br><p>eBay Selling and online Sales Training. Basic - Advanced - Sign up for Bids4U's eBay Selling Tips Newsletter Email: eBay.com eBay Training Books Free eBay! Registration eBay Radio eBay Bible eBay ... </P></li><li><A href="http://pages.ebay.com/education/" target='_blank'>eBay: Learning Center: Index</A><br><p>New to eBay? Welcome! This is the place to start. These courses will give you a quick overview of how eBay works and how to buy with confidence. This symbol indicates an Audio Tour, please turn up ... </P></li></ul>
						<br>
						<br>
						<Center><font color='CCCCCC'>Back</font> | 1 | <a href='Ebay_Training_Page_2.php'>2</a> | <a href='Ebay_Training_Page_3.php'>3</a> | <a href='Ebay_Training_Page_2.php'>Next</a><br><br></center>
					</p>
					<br>
					<?php if ($DisplayAmazon) { echo "<br><center>"; 
					include("../Includes/Amazon_468x60.php"); echo "</center><br>"; } ?>
				</td>
				<td class="tdRight">
					<?php include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="3" class="tdRow">
					<?=$Category;?> News and Current Events
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include($Domain."Includes/NewsFeed.php?CAT=$NewsFeed"); ?> <?php if 
					($DisplayAmazon) { echo "<br><center>"; 
					include("../Includes/Amazon_728x90.php"); echo "</center><br>"; } ?>
				</td>
			</tr>
			<tr>
				<td class="tdFooter" colspan="3">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All 
					Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal 
						Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>

