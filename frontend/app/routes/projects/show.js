import Ember from 'ember';

export default Ember.Route.extend({
  model: function(params) {
    return this.store.find('project', params.id);
  },

  afterModel: function(model, transition) {
    if (transition.intent.url == null) { // Prevents double loading the model.
      model.reload()
    }
  }
});
