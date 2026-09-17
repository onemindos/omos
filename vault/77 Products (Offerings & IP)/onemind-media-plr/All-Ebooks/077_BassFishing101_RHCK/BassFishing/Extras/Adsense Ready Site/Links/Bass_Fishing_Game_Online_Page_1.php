<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title>Links | Bass Fishing Game Online Page 1</title>
		<META http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Bass Fishing Game Online Related Links - Page 1 - ">
		<meta name="Keywords"			content="<?=$Keywords;?>">
		<meta name="Distribution"		content="global">
		<meta name="Publisher"			content="<?=$Domain;?>">
		<meta name="Rating"				content="General">
		<meta name="Revisit-after"		content="5 days">
		<meta name="Robots"				content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0">
			<tr>
				<td class="tdHeader" colspan="2">
					<p style="margin-bottom: 5; margin-right: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td class="tdContent">
					<br>
					<p>
						<b>Bass Fishing Game Online Related Links - Page 1 </b>
						<br>
						<b>
							<small>
								
							</small>
						</b>
						<br>
						<br>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<ul class='ulCategory'><li><A href="http://www.stripersonline.com/" target='_blank'>Striped Bass Fishing - Stripers Online - Striper Fishing Online!</A><br><p>Tips for catching more fish. Includes articles, reports, chats and forums for the New Jersey shores and an online store.</P></li><li><A href="http://www.arkadium.com/fishing.html" target='_blank'>Flash Fishing Games, Online Fishing Game, Internet Fishing Games, Web ... </A><br><p>Your source for online fishing games, internet fishing games, web fishing games and bass fishing game online Fishing games at Arkadium. Your source for fishing games, flash fishing game and fishing ... </P></li></ul>
						<br>
						<br>
						<Center><font color='CCCCCC'>Back</font> | 1 | <a href='Bass_Fishing_Game_Online_Page_2.php'>2</a> | <a href='Bass_Fishing_Game_Online_Page_3.php'>3</a> | <a href='Bass_Fishing_Game_Online_Page_2.php'>Next</a><br><br></center>
					</p>
					<br>
					<?php	
						if ($DisplayAmazon)
							{
								echo "<br><center>";
								include("../Includes/Amazon_468x60.php");
								echo "</center><br>";
							}
					?>					
				</td>
				<td class="tdRight">
					<?php include("../Includes/Navigation.php"); ?>
					<?php include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="2">
					<?=$Category;?> News and Events
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<br>
					<?php include($Domain."Includes/NewsFeed.php?CAT=Bass+Fishing+Game+Online"); ?>
					<?php	
						if ($DisplayAmazon)
							{
								echo "<br><center>";
								include("../Includes/Amazon_728x90.php");
								echo "</center><br>";
							}
					?>
				</td>
			</tr>
			<tr>
				<td class="tdFooter" colspan="2">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>

