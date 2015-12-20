import DS from 'ember-data';

export default DS.RESTSerializer.extend({
  keyForRelationship: function(key, relationship) {
    if (key == "branches") {
      return "branch_names";
    } else {
      return key;
    }
  }
});
