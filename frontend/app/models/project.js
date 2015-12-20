import DS from 'ember-data';

export default DS.Model.extend({
  name: DS.attr('string'),
  repo: DS.attr('string'),
  branches: DS.hasMany('branch', { async: true })
});
