<?php include("Session.php"); ?>
<?php

	$cfgData = file_get_contents("Configure.tpl");

    $cfgData = str_replace("%UserName%",		$_SESSION[UserName],		$cfgData);
    $cfgData = str_replace("%Password%",		$_SESSION[Password],		$cfgData);
    $cfgData = str_replace("%SiteName%",		$_SESSION[SiteName],		$cfgData);
    $cfgData = str_replace("%Category%",		$_SESSION[Category],		$cfgData);
    $cfgData = str_replace("%NewsFeed%",		$_SESSION[NewsFeed],		$cfgData);
    $cfgData = str_replace("%MainTitle%",		$_SESSION[MainTitle],		$cfgData);
    $cfgData = str_replace("%SubTitle%",		$_SESSION[SubTitle],		$cfgData);
    $cfgData = str_replace("%Keywords%",		$_SESSION[Keywords],		$cfgData);
    $cfgData = str_replace("%GooglePubID%",		$_SESSION[GooglePubID],		$cfgData);
    $cfgData = str_replace("%GoogleChannel%",	$_SESSION[GoogleChannel],	$cfgData);
    $cfgData = str_replace("%GoogleLinkColor%",	$_SESSION[GoogleLinkColor],	$cfgData);
    $cfgData = str_replace("%GoogleURLColor%",	$_SESSION[GoogleURLColor],	$cfgData);
    $cfgData = str_replace("%Name%",			$_SESSION[Name],			$cfgData);
    $cfgData = str_replace("%Address1%",		$_SESSION[Address1],		$cfgData);
    $cfgData = str_replace("%Address2%",		$_SESSION[Address2],		$cfgData);
    $cfgData = str_replace("%City%",			$_SESSION[City],			$cfgData);
    $cfgData = str_replace("%State%",			$_SESSION[State],			$cfgData);
    $cfgData = str_replace("%ZipCode%",			$_SESSION[ZipCode],			$cfgData);
    $cfgData = str_replace("%Country%",			$_SESSION[Country],			$cfgData);
    $cfgData = str_replace("%Telephone%",		$_SESSION[Telephone],		$cfgData);
    $cfgData = str_replace("%EmailAddress%",	$_SESSION[EmailAddress],	$cfgData);
    $cfgData = str_replace("%ShowGoogle%",		$_SESSION[ShowGoogle],		$cfgData);
    $cfgData = str_replace("%ShowFireFox%",		$_SESSION[ShowFireFox],		$cfgData);	
    $cfgData = str_replace("%ShowSearch%",		$_SESSION[ShowSearch],		$cfgData);	

	if (isset($_SESSION['UserName']))
		{
			$fp = fopen("../Includes/Configure.php", "w");
			$fw = fwrite( $fp, $cfgData );
			fclose( $fp );
		}	
?>
	<?php include("Header.php"); ?>
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							SYSTEM MESSAGE&nbsp; 
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<small>
								&nbsp;
							</small>
						</td>				
					<tr>
					<tr>
						<td>
							Your Niche Mania Website has been Saved with the new settings you entered!
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldValue">
							<input style="width: 50%;" onclick="window.open('../index.php');" type="button" value=" Preview Website ">					
						</td>				
					<tr>
				</table>
			</center>
	<?php include("Footer.php"); ?>