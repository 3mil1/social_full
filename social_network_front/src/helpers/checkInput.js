export const handleInputs = (elementName, input) => {
  if (!input || input == ' ' || input == '') {
    document.getElementById(elementName).value = 'fill ';
    document.getElementById(elementName).classList.add('error');
    return false;
  }
  return true;
};

 export const handleAfterErrorClick = elementName => {
   if (document.getElementById(elementName).classList.contains('error')) {
     document.getElementById(elementName).classList.remove('error');
     document.getElementById(elementName).value = '';
   }
 };

export default handleInputs;