					<tr>
						<td colspan="2">
							<?php include("Error.php"); ?>
						</td>				
					</tr>
					<tr>
						<td class="tdFieldName">
							Keyword
						</td>									
						<td>
							<select name="Keywords" style="width: 100%;">
								<?php
									$Array			= explode(",", $_SESSION['Keywords']);
									$intCounter		= 0;
									
									while($intCounter != count($Array))
										{
											if(trim($Keyword) == trim($Array[$intCounter]))
												{
													echo "<option selected value='".$Array[$intCounter]."'>".$Array[$intCounter]."</option>";
												}
											else
												{
													echo "<option value='".$Array[$intCounter]."'>".$Array[$intCounter]."</option>";
												}
											$intCounter++;
										}									
								?>
								
							</select>	
						</td>				
					</tr>
					<tr>
						<td class="tdFieldName">
							Article Title
						</td>															
						<td align="center">
							<input class="iText" type="text" maxsize="200" style="width: 100%;" name="Title" value="<?=$Title;?>">
						</td>				
					</tr>
					<tr>
						<td class="tdFieldName">
							Author's Name
						</td>															
						<td align="center">
							<input class="iText" maxsize="200" type="text" style="width: 100%;" name="Author" value="<?=stripslashes($Author);?>">
						</td>				
					</tr>
					<tr>
						<td colspan="2">
							<textarea id="HTML_Article" name="Article"><?=stripslashes($Article);?></textarea>
						</td>				
					<tr>
					<tr>
						<td colspan="2">
							<textarea id="HTML_Resource" name="Resource"><?=stripslashes($Resource);?></textarea>
						</td>				
					</tr>
					<tr>
						<td colspan="2">
							<hr size="0">
						</td>				
					</tr>
					<tr>
						<td class="tdFieldValue" colspan="2">
							<input class="iButtonEx" onClick="updateTextArea('HTML_Resource'); updateTextArea('HTML_Article'); this.disabled='true'; frmData.submit();" type="submit" value=" Save Article ">
							<input class="iButtonEx" onClick="this.disabled='true'; location.href='Articles.php'" type="button" value=" Cancel ">
						</td>	
					</tr>
				</table>
			</center>
	<?php include("Footer.php"); ?>
	<script language="javascript1.2">
		generate_wysiwyg('HTML_Article', 635, 500);
		generate_wysiwyg('HTML_Resource', 635, 150);
	</script>