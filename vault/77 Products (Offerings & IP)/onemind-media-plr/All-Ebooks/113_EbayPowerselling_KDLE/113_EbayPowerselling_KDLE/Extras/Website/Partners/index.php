<?php include("../Includes/Configure.php"); ?>
<html>
	<head>
		<title><?=$Category;?> Partners</title>
		<META http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Kotzakoliou, SSA - http://www.niche-mania.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="<?=$Category;?> Partners">
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
				<td class="tdHeader" colspan="3">
					<p style="margin-bottom: 5; margin-left: 10;">
						<?php include("../Includes/Menu.php"); ?>
					</p>
				</td>
			</tr>
			<tr>
				<td colspan="32">
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
					<p>
						<b>
							<?=$SiteName;?> Partners
						</b>
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>					
					<?php include("../Includes/Partners.php"); ?>
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
					<?php include($Domain."Includes/NewsFeed.php?CAT=$NewsFeed"); ?>
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
				<td class="tdFooter" colspan="3">
					&copy; 2005, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>
