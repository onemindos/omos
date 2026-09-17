<?php include("Session.php"); ?>
<?php

	$strAmazon								= "disabled";
	
	if (!$_SESSION['AmazonLinkColor'])
		{
			$_SESSION['AmazonLinkColor']	= $_SESSION['GoogleLinkColor'];
			$_SESSION['AmazonTextColor']	= $_SESSION['GoogleURLColor'];
		}
		
	if (isset($_POST['AmazonID']))
		{						
			$_SESSION['AmazonID']			= $_POST['AmazonID'];
			$_SESSION['AmazonMode']			= $_POST['AmazonMode'];
			$_SESSION['AmazonSearch']		= $_POST['AmazonSearch'];
			$_SESSION['AmazonLinkColor']	= $_POST['AmazonLinkColor'];
			$_SESSION['AmazonTextColor']	= $_POST['AmazonTextColor'];
			$_SESSION['DisplayAmazon']		= $_POST['DisplayAmazon'];
		
			$Amazon_728x90					= "<iframe src='http://rcm.amazon.com/e/cm?t=$AmazonID&o=1&p=48&l=st1&mode=$AmazonMode&search=$AmazonSearch&fc1=$AmazonTextColor&lt1=_blank&lc1=$AmazonLinkColor&bg1=&f=ifr' marginwidth='0' marginheight='0' width='728' height='90' border='0' frameborder='0' style='border:none;' scrolling='no'></iframe>";
			$Amazon_468x60					= "<iframe src='http://rcm.amazon.com/e/cm?t=$AmazonID&o=1&p=13&l=st1&mode=$AmazonMode&search=$AmazonSearch&fc1=$AmazonTextColor&lt1=_blank&lc1=$AmazonLinkColor&bg1=&f=ifr' marginwidth='0' marginheight='0' width='468' height='60' border='0' frameborder='0' style='border:none;' scrolling='no'></iframe>";
			
			
			$fp								= fopen("../Includes/Amazon_728x90.php", "w");
			$fw								= fwrite( $fp, $Amazon_728x90);
			fclose($fp);

			$fp								= fopen("../Includes/Amazon_468x60.php", "w");
			$fw								= fwrite( $fp, $Amazon_468x60);
			fclose($fp);
		
			header('location: Write.php');
		}
?>
	<?php include("Header.php"); ?>	
		<form method="post" action="Amazon.php" name="frmData">
			<center>
				<table class="tblForm">
					<tr>
						<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;" colspan="2">
							Amazon Settings 
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							&nbsp;
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Associates ID
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="AmazonID" value="<?=$_SESSION['AmazonID']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Category
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="AmazonMode" value="<?=$_SESSION['AmazonMode']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Keyword
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="AmazonSearch" value="<?=$_SESSION['AmazonSearch']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Link Color
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="AmazonLinkColor" value="<?=$_SESSION['AmazonLinkColor']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Text Color
						</td>				
						<td class="tdFieldValue">
							<input class="iText" type="text" style="width: 100%;" name="AmazonTextColor" value="<?=$_SESSION['AmazonTextColor']; ?>">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
							Display Amazon Ads
						</td>				
						<td align="left">
							<input type="checkbox" Name="DisplayAmazon" value="checked" <?=$_SESSION['DisplayAmazon'];?>>
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					<tr>
					<tr>
						<td class="tdFieldName">
						</td>				
						<td class="tdFieldValue">
							<input class="iButton" onClick="this.disabled='true'; frmData.submit();" type="submit" value=" Save Amazon Settings ">					
						</td>				
					<tr>
				</table>
			</center>
		</form>
	<?php include("Footer.php"); ?>