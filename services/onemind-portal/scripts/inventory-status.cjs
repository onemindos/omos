// Absence in discovery is not proof an application or device is offline/deleted.
const missingPrefix = 'Not observed in the latest inventory; review required. ';
function missingDescription(description) {
  const previous = typeof description === 'string' ? description : 'No previous description recorded.';
  return previous.startsWith(missingPrefix) ? previous : missingPrefix + previous;
}
module.exports = { missingDescription };
