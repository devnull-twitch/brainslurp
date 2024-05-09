(() => {
  const listEntryTpl = document.querySelector("#nested-list-item");
  const listTagTpl = document.querySelector("#nested-list-tag");

  // utility method to create 
  const addTagWithLabel = (box, colorClass, tagCategory, fieldID) => {
    const orgField = document.querySelector(fieldID);
    if (orgField.value.replace(' ', '').length <= 0) {
      return;
    }

    const containerSpan = document.createElement('span');
    containerSpan.classList.add('mr-1');

    const prefix = document.createElement('span');
    prefix.innerText = tagCategory + '( ';
    containerSpan.appendChild(prefix);
    
    const tag = listTagTpl.content.cloneNode(true).firstChild;
    if (orgField.tagName && orgField.tagName.toLowerCase() == 'select') {
      tag.innerText = orgField.querySelector(`option[value="${orgField.value}"]`).innerText;
    } else {
      tag.innerText = orgField.value;
    }
    tag.classList.add(colorClass);
    containerSpan.appendChild(tag);

    const postfix = document.createElement('span');
    postfix.innerText = ' )';
    containerSpan.appendChild(postfix);

    box.appendChild(containerSpan);
  }

  const persistField = (box, orgID, targetName) => {
    const orgField = document.querySelector(orgID);
    const hiddenField = document.createElement('input');
    hiddenField.name = targetName;
    hiddenField.type = 'hidden';
    hiddenField.value = orgField.value;
    orgField.value = '';
    box.appendChild(hiddenField);
  };

  // Click handler for requirements
  document.querySelector('.brainslurp-add-requirement button').addEventListener('click', () => {
    const reqList = document.querySelector('#req-box ol');
    const newListEntry = listEntryTpl.content.cloneNode(true).firstChild;

    addTagWithLabel(newListEntry, 'bg-lime-500', 'with', '#tpl_req_tags');
    addTagWithLabel(newListEntry, 'bg-red-500', 'without', '#tpl_req_no_tags');
    addTagWithLabel(newListEntry, 'bg-sky-500', 'category', '#tpl_req_category');

    persistField(newListEntry, '#tpl_req_tags', 'req_tags');
    persistField(newListEntry, '#tpl_req_no_tags', 'req_no_tags');
    persistField(newListEntry, '#tpl_req_category', 'req_category');

    // delete placeholder "None" element if it is still around
    const placeHolder = reqList.firstChild;
    if (placeHolder.classList.contains('placeholder')) {
      reqList.removeChild(placeHolder);
    }

    reqList.appendChild(newListEntry);
  });

  // Click handler for actions
  document.querySelector('.brainslurp-add-action').addEventListener('click', () => {
    const actionList = document.querySelector('#action-box ol');
    const newListEntry = listEntryTpl.content.cloneNode(true).firstChild;

    addTagWithLabel(newListEntry, 'bg-sky-500', 'name', '#tpl_action_name');
    addTagWithLabel(newListEntry, 'bg-lime-500', 'add', '#tpl_action_adds');
    addTagWithLabel(newListEntry, 'bg-red-500', 'remove', '#tpl_action_removes');

    persistField(newListEntry, '#tpl_action_name', 'action_name');
    persistField(newListEntry, '#tpl_action_adds', 'action_adds');
    persistField(newListEntry, '#tpl_action_removes', 'action_removes');

    // delete placeholder "None" element if it is still around
    const placeHolder = actionList.firstChild;
    if (placeHolder.classList.contains('placeholder')) {
      actionList.removeChild(placeHolder);
    }

    actionList.appendChild(newListEntry);
  });
})();